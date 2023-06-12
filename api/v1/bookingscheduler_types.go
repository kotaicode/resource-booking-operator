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

// BookingSchedulerSpec defines the desired state of BookingScheduler
type BookingSchedulerSpec struct {
	Schedule        string      `json:"schedule,omitempty"`
	Duration        int         `json:"duration,omitempty"`
	BookingTemplate BookingSpec `json:"bookingTemplate,omitempty"`
}

// BookingSchedulerStatus defines the observed state of BookingScheduler
type BookingSchedulerStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// BookingScheduler is the Schema for the bookingschedulers API
type BookingScheduler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BookingSchedulerSpec   `json:"spec,omitempty"`
	Status BookingSchedulerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BookingSchedulerList contains a list of BookingScheduler
type BookingSchedulerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BookingScheduler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BookingScheduler{}, &BookingSchedulerList{})
}
