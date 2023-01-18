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

	managerv1 "github.com/kotaicode/resource-booking-operator/api/v1"
	"github.com/kotaicode/resource-booking-operator/clients"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ResourceMonitorReconciler reconciles a ResourceMonitor object
type ResourceMonitorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const (
	ResourceType      = "ec2"
	ResourceNamespace = "default"
)

//+kubebuilder:rbac:groups=manager.kotaico.de,resources=resourcemonitors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=manager.kotaico.de,resources=resourcemonitors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=manager.kotaico.de,resources=resourcemonitors/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ResourceMonitor object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile

func getDistinctTags(uniqueTags []string, clusterResources []string) []string {
	var distinctTags []string

	for i := 0; i < 2; i++ {
		for _, s1 := range uniqueTags {
			found := false
			for _, s2 := range clusterResources {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				distinctTags = append(distinctTags, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			uniqueTags, clusterResources = clusterResources, uniqueTags
		}
	}

	return distinctTags
}

func (r *ResourceMonitorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	var clusterResources []string

	log.Info("Reconcile resource monitor")

	var resourceMonitor managerv1.ResourceMonitor
	if err := r.Get(ctx, req.NamespacedName, &resourceMonitor); err != nil {
		log.Error(err, "Error getting resource monitor")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var resources managerv1.ResourceList
	if err := r.List(context.Background(), &resources); err != nil {
		log.Error(err, "Error listing resource monitor")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	for _, rs := range resources.Items {
		clusterResources = append(clusterResources, rs.Spec.Tag)
	}

	uniqueTags, err := clients.GetUniqueTags()
	if err != nil {
		//TODO: condition not satisfied ?
		log.Info("getting unique tags success")
	}
	nonMatchingTags := getDistinctTags(uniqueTags, clusterResources)

	for _, tag := range nonMatchingTags {
		resource := &managerv1.Resource{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "manager.kotaico.de/v1",
				Kind:       "Resource",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      tag,
				Namespace: ResourceNamespace,
			},
			Spec: managerv1.ResourceSpec{
				Booked: false,
				Tag:    tag,
				Type:   ResourceType,
			},
		}
		log.Info("creating resources")
		r.Create(ctx, resource)
	}

	return ctrl.Result{RequeueAfter: time.Duration(time.Second * 15)}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ResourceMonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&managerv1.ResourceMonitor{}).
		Complete(r)
}
