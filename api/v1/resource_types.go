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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ResourceSpec defines the desired state of Resource
type ResourceSpec struct {
	BookedBy    string `json:"booked_by"`
	BookedUntil string `json:"booked_until"`

	Tag  string `json:"tag"`
	Type string `json:"type"`
}

// ResourceStatus defines the observed state of Resource
type ResourceStatus struct {
	Instances   int    `json:"instances"`
	Running     int    `json:"running"`
	Status      string `json:"status"`
	LockedBy    string `json:"locked_by"`
	LockedUntil string `json:"locked_until"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:JSONPath=".status.locked_by",name="LOCKED BY",type="string"
//+kubebuilder:printcolumn:JSONPath=".status.locked_until",name="LOCKED UNTIL",type="string"
//+kubebuilder:printcolumn:JSONPath=".status.instances",name="INSTANCES",type="integer"
//+kubebuilder:printcolumn:JSONPath=".status.running",name="RUNNING",type="integer"
//+kubebuilder:printcolumn:JSONPath=".status.status",name="STATUS",type="string"

// Resource is the Schema for the resources API
type Resource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceSpec   `json:"spec,omitempty"`
	Status ResourceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ResourceList contains a list of Resource
type ResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Resource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Resource{}, &ResourceList{})
}
