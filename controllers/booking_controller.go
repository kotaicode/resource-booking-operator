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
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	managerv1 "github.com/kotaicode/resource-booking-operator/api/v1"
	"github.com/kotaicode/resource-booking-operator/notify"
)

// BookingReconciler reconciles a Booking object
type BookingReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=manager.kotaico.de,resources=bookings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=manager.kotaico.de,resources=bookings/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=manager.kotaico.de,resources=bookings/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *BookingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var err error
	log := log.FromContext(ctx)
	log.Info("Reconciling resource")

	var booking managerv1.Booking
	if err := r.Get(ctx, req.NamespacedName, &booking); err != nil {
		log.Error(err, "Error getting booking")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	bookStart, err := time.Parse(time.RFC3339, booking.Spec.StartAt)
	if err != nil {
		log.Error(err, "Error parsing booking start")
	}

	bookEnd, err := time.Parse(time.RFC3339, booking.Spec.EndAt)
	if err != nil {
		log.Error(err, "Error parsing booking end")
	}

	resNamespacedName := types.NamespacedName{
		Namespace: req.Namespace,
		Name:      booking.Spec.ResourceName,
	}

	var resource managerv1.Resource
	if err := r.Get(context.Background(), resNamespacedName, &resource); err != nil {
		log.Error(err, "Error listing bookings")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if bookStart.Before(time.Now()) && time.Now().Before(bookEnd) {
		booking.Status.Status = managerv1.BookingInProgress
		updateResource(r, ctx, &resource, &booking)
	} else if bookEnd.Before(time.Now()) {
		booking.Status.Status = managerv1.BookingFinished
		updateResource(r, ctx, &resource, &booking)
	} else {
		booking.Status.Status = managerv1.BookingScheduled
	}

	if booking.Status.Status == managerv1.BookingInProgress && time.Until(bookEnd) < time.Minute*20 &&
		!booking.Status.NotificationSent && len(booking.Spec.Notifications) > 0 {

		for _, notification := range booking.Spec.Notifications {
			n, err := notify.NewNotifier(notification)
			if err != nil {
				log.Error(err, "Error sending notification")
			}

			err = n.Prepare(booking).Send()
			if err != nil {
				log.Error(err, "Error sending notification")
			}

			booking.Status.NotificationSent = err == nil
		}
	}

	log.Info("Updating booking status", "status", booking.Status.Status)
	err = r.Status().Update(ctx, &booking)
	if err != nil {
		log.Error(err, "Error updating booking status")
		return ctrl.Result{}, err
	}

	if booking.Status.Status == managerv1.BookingFinished {
		log.Info("Booking finished")
		return ctrl.Result{}, nil
	}

	return ctrl.Result{RequeueAfter: time.Duration(time.Minute * 1)}, nil
}

func updateResource(r *BookingReconciler, ctx context.Context, rs *managerv1.Resource, booking *managerv1.Booking) {
	log := log.FromContext(ctx)

	if booking.Status.Status == managerv1.BookingInProgress {
		rs.Spec.BookedBy = booking.Spec.UserID
		rs.Spec.BookedUntil = booking.Spec.EndAt
	} else if booking.Status.Status == managerv1.BookingFinished {
		rs.Spec.BookedUntil = ""
		rs.Spec.BookedBy = ""
	}

	err := r.Update(ctx, rs)
	if err != nil {
		log.Error(err, "Error updating resource spec")
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *BookingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&managerv1.Booking{}).
		Complete(r)
}
