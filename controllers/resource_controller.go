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
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	managerv1 "my.domain/resource-booking/api/v1"
	"my.domain/resource-booking/instances"
)

// ResourceReconciler reconciles a Resource object
type ResourceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=manager.my.domain,resources=resources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=manager.my.domain,resources=resources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=manager.my.domain,resources=resources/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Resource object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *ResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var resource managerv1.Resource
	if err := r.Get(ctx, req.NamespacedName, &resource); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	res := instances.Resource{NameTag: resource.Spec.Tag}
	instanceStatus, _ := res.Status()

	var running, available int32
	var status string
	for _, v := range instanceStatus {
		available++
		if v.InstanceStatusCode == instances.StatusRunning {
			running++
		}
	}

	if running == 0 {
		status = "STOPPED"
	} else if running == available {
		status = "RUNNING"
	} else {
		status = "PENDING"
	}

	if resource.Spec.Booked {
		// Just so I can test
		if status == "STOPPED" {
			res.Start()
		}
	} else {
		// Just so I can test
		if status == "RUNNING" {
			res.Stop()
		}
	}

	resource.Status = managerv1.ResourceStatus{
		Instances: available,
		Running:   running,
		Status:    status,
	}

	err := r.Status().Update(ctx, &resource)
	if err != nil {
		return ctrl.Result{}, err
	}

	// log.Info("reconciled resource")
	return ctrl.Result{RequeueAfter: time.Duration(time.Second * 30)}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&managerv1.Resource{}).
		Complete(r)
}
