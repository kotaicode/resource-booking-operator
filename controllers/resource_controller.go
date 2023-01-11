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

	managerv1 "github.com/kotaicode/resource-booking-operator/api/v1"
	"github.com/kotaicode/resource-booking-operator/clients"
)

// ResourceReconciler reconciles a Resource object
type ResourceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=manager.kotaico.de,resources=resources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=manager.kotaico.de,resources=resources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=manager.kotaico.de,resources=resources/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *ResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var status string

	log := log.FromContext(ctx)
	log.Info("Reconciling resource")

	var resource managerv1.Resource
	if err := r.Get(ctx, req.NamespacedName, &resource); err != nil {
		log.Error(err, "Error getting resource")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	cloudResource, err := clients.ResourceFactory(resource.Spec.Type, resource.Spec.Tag)
	if err != nil {
		log.Error(err, err.Error())
		return ctrl.Result{}, err
	}

	rStat, err := cloudResource.Status()
	if err != nil {
		log.Error(err, "Error getting resource status")
		return ctrl.Result{}, err
	}

	switch rStat.Running {
	case 0:
		status = clients.StatusStopped
	case rStat.Available:
		status = clients.StatusRunning
	default:
		status = clients.StatusPending
	}

	if resource.Spec.Booked {
		if status != clients.StatusRunning {
			if err := cloudResource.Start(resource.Status.LockedBy, resource.Status.LockedUntil); err != nil {
				log.Error(err, "Error starting resource instances")
				return ctrl.Result{RequeueAfter: time.Duration(time.Second * 60)}, err
			}
		}
	} else {
		if status == clients.StatusRunning {
			if err := cloudResource.Stop(resource.Status.LockedBy); err != nil {
				log.Error(err, "Error stopping resource instances")
				return ctrl.Result{RequeueAfter: time.Duration(time.Second * 60)}, err
			}
		}
	}

	resource.Status = managerv1.ResourceStatus{
		LockedBy:    string(resource.ObjectMeta.UID),
		LockedUntil: resource.Status.LockedUntil,
		Instances:   rStat.Available,
		Running:     rStat.Running,
		Status:      status,
	}

	err = r.Status().Update(ctx, &resource)
	if err != nil {
		log.Error(err, "Error updating resource status")
		return ctrl.Result{}, err
	}

	log.Info("reconciled resource")
	return ctrl.Result{RequeueAfter: time.Duration(time.Second * 15)}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&managerv1.Resource{}).
		Complete(r)
}
