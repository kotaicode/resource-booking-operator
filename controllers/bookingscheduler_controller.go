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
		log.Error(err, "Error parsing schedule", "schedule", bookingScheduler.Spec.Schedule)
		return ctrl.Result{}, err
	}

	now := time.Now()

	nextSched := schedule.Next(now)
	nextInMin := nextSched.Sub(now)

	booking = setBooking(bookingScheduler, booking)

	if bookingScheduler.Status.Next != "" {
		statusNext, err := time.Parse(time.RFC3339, bookingScheduler.Status.Next)
		if err != nil {
			log.Error(err, "Error parsing status.next", "value", bookingScheduler.Status.Next)
			return ctrl.Result{}, err
		}

		if statusNext.Equal(now) || statusNext.Before(now) {
			if err := r.Create(ctx, &booking); err != nil {
				log.Error(err, "Error creating booking")
				return ctrl.Result{}, err
			}
		}
	}

	bookingScheduler.Status.Next = nextSched.Format(time.RFC3339)

	err = r.Status().Update(ctx, &bookingScheduler)
	if err != nil {
		log.Error(err, "Error updating booking scheduler status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Duration(nextInMin)}, nil
}

// setBooking grabs the necessary information from the booking scheduler and sets it to the booking
func setBooking(bookingScheduler managerv1.BookingScheduler, booking managerv1.Booking) managerv1.Booking {
	booking.Spec = bookingScheduler.Spec.BookingTemplate

	now := time.Now()
	booking.Spec.StartAt = now.Format(time.RFC3339)

	endAt := now.Add(time.Duration(bookingScheduler.Spec.Duration) * time.Minute)
	booking.Spec.EndAt = endAt.Format(time.RFC3339)

	booking.Name = fmt.Sprintf("%s-%s-%d", booking.Spec.UserID, booking.Spec.ResourceName, endAt.Unix())
	booking.Namespace = bookingScheduler.Namespace

	return booking
}

// SetupWithManager sets up the controller with the Manager.
func (r *BookingSchedulerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&managerv1.BookingScheduler{}).
		Complete(r)
}
