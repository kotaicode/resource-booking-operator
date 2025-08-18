/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	managerv1 "github.com/kotaicode/resource-booking-operator/api/v1"
)

// BookingSchedulerReconciler reconciles a BookingScheduler object
type BookingSchedulerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=manager.kotaico.de,resources=bookingschedulers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=manager.kotaico.de,resources=bookingschedulers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=manager.kotaico.de,resources=bookingschedulers/finalizers,verbs=update

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *BookingSchedulerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling BookingScheduler", "name", req.Name, "namespace", req.Namespace)

	var (
		bookingScheduler managerv1.BookingScheduler
		booking          managerv1.Booking
	)

	if err := r.Get(ctx, req.NamespacedName, &bookingScheduler); err != nil {
		log.Error(err, "Error getting booking scheduler")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	schedule, err := cron.ParseStandard(bookingScheduler.Spec.Schedule)
	if err != nil {
		log.Error(err, "Error parsing schedule")
		return ctrl.Result{}, err
	}

	now := time.Now().UTC()
	log.Info("Processing scheduler", "schedule", bookingScheduler.Spec.Schedule, "currentTime", now.Format("15:04"), "statusNext", bookingScheduler.Status.Next)

	nextSched := schedule.Next(now)
	nextInMin := nextSched.Sub(now)
	log.Info("Next scheduled run", "nextRun", nextSched.Format("2006-01-02 15:04"), "minutesUntilNext", int(nextInMin.Minutes()))

	booking = setBooking(bookingScheduler, booking)
	log.Info("Prepared booking", "bookingName", booking.Name, "startAt", booking.Spec.StartAt, "endAt", booking.Spec.EndAt)

	// Check if we should create a booking now
	shouldCreateNow := false

	if bookingScheduler.Status.Next == "" {
		// First run - always create a booking for the current day if we have start/end times
		shouldCreateNow = true
	} else {
		// Check if the scheduled time has arrived
		statusNext, err := time.Parse(time.RFC3339, bookingScheduler.Status.Next)
		if err != nil {
			log.Error(err, "Error parsing status.next")
			return ctrl.Result{}, err
		}

		if statusNext.Equal(now) || statusNext.Before(now) {
			shouldCreateNow = true
		}
	}

	// Create the booking if it's time
	if shouldCreateNow {
		log.Info("Creating booking", "scheduler", bookingScheduler.Name, "booking", booking.Name)
		if err := r.Create(ctx, &booking); err != nil {
			log.Error(err, "Error creating booking")
			return ctrl.Result{}, err
		}
		log.Info("Successfully created booking", "scheduler", bookingScheduler.Name, "booking", booking.Name)
	}

	bookingScheduler.Status.Next = nextSched.Format(time.RFC3339)

	err = r.Status().Update(ctx, &bookingScheduler)
	if err != nil {
		log.Error(err, "Error updating booking scheduler status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Duration(nextInMin)}, nil
}

// setBooking creates a booking based on the scheduler's configuration
func setBooking(bookingScheduler managerv1.BookingScheduler, booking managerv1.Booking) managerv1.Booking {
	now := time.Now().UTC()

	// Check if we should use time-based scheduling (start_time/end_time) or timestamp-based (start_at/end_at)
	if bookingScheduler.Spec.StartTime != "" && bookingScheduler.Spec.EndTime != "" {
		// Use time-based scheduling - create booking for today with specified start/end times
		startTime, err := time.Parse("15:04", bookingScheduler.Spec.StartTime)
		if err != nil {
			// If parsing fails, use current time as fallback
			startTime = now
		}

		endTime, err := time.Parse("15:04", bookingScheduler.Spec.EndTime)
		if err != nil {
			// If parsing fails, use current time + 1 hour as fallback
			endTime = now.Add(time.Hour)
		}

		// Create the actual start and end times for today
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		actualStartTime := today.Add(time.Duration(startTime.Hour())*time.Hour + time.Duration(startTime.Minute())*time.Minute)
		actualEndTime := today.Add(time.Duration(endTime.Hour())*time.Hour + time.Duration(endTime.Minute())*time.Minute)

		// Set the booking details using calculated times
		booking.Spec.StartAt = actualStartTime.Format(time.RFC3339)
		booking.Spec.EndAt = actualEndTime.Format(time.RFC3339)
	} else if bookingScheduler.Spec.StartAt != "" && bookingScheduler.Spec.EndAt != "" {
		// Use existing timestamp-based scheduling
		booking.Spec.StartAt = bookingScheduler.Spec.StartAt
		booking.Spec.EndAt = bookingScheduler.Spec.EndAt
	} else {
		// Fallback: use current time + 1 hour
		booking.Spec.StartAt = now.Format(time.RFC3339)
		booking.Spec.EndAt = now.Add(time.Hour).Format(time.RFC3339)
	}

	// Set the other booking details
	booking.Spec.ResourceName = bookingScheduler.Spec.ResourceName
	booking.Spec.UserID = bookingScheduler.Spec.UserID
	booking.Spec.Notifications = bookingScheduler.Spec.Notifications

	// Generate a unique name for the booking
	startTime, _ := time.Parse(time.RFC3339, booking.Spec.StartAt)
	booking.Name = fmt.Sprintf("%s-%s-%d", booking.Spec.UserID, booking.Spec.ResourceName, startTime.Unix())
	booking.Namespace = bookingScheduler.Namespace

	return booking
}

// SetupWithManager sets up the controller with the Manager.
func (r *BookingSchedulerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&managerv1.BookingScheduler{}).
		Complete(r)
}
