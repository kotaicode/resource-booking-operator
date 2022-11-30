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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	managerv1 "github.com/kotaicode/resource-booking-operator/api/v1"
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
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Booking object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *BookingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var resources managerv1.ResourceList
	var booking managerv1.Booking
	if err := r.Get(ctx, req.NamespacedName, &booking); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	bookStart, err := time.Parse(time.RFC3339, booking.Spec.StartAt)
	if err != nil {
		fmt.Println("TODO ERROR", err.Error())
	}

	bookEnd, err := time.Parse(time.RFC3339, booking.Spec.EndAt)
	if err != nil {
		fmt.Println("TODO ERROR", err.Error())
	}

	if err := r.List(context.Background(), &resources, client.MatchingFields{"spec.tag": booking.Spec.ResourceName}); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if bookStart.Before(time.Now()) && time.Now().Before(bookEnd) {
		updateResource(r, ctx, &resources, true)
		booking.Status.Status = managerv1.BookingInProgress
	} else if bookStart.Before(time.Now()) && bookEnd.Before(time.Now()) {
		updateResource(r, ctx, &resources, false)
		booking.Status.Status = managerv1.BookingFinished
	} else {
		updateResource(r, ctx, &resources, false)
		booking.Status.Status = managerv1.BookingScheduled
	}

	err = r.Status().Update(ctx, &booking)
	if err != nil {
		return ctrl.Result{}, err
	}

	if booking.Status.Status == managerv1.BookingFinished {
		return ctrl.Result{}, nil
	}

	return ctrl.Result{RequeueAfter: time.Duration(time.Minute * 1)}, nil
}

func updateResource(r *BookingReconciler, ctx context.Context, resources *managerv1.ResourceList, booked bool) {
	for _, rs := range resources.Items {
		rs.Spec.Booked = booked

		err := r.Update(ctx, &rs)
		if err != nil {
			fmt.Println("TODO")
		}
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *BookingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	mgr.GetFieldIndexer().IndexField(context.TODO(), &managerv1.Resource{}, "spec.tag", func(o client.Object) []string {
		return []string{o.(*managerv1.Resource).Spec.Tag}
	})

	return ctrl.NewControllerManagedBy(mgr).
		For(&managerv1.Booking{}).
		Complete(r)
}