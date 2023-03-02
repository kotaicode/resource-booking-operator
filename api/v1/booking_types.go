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

const (
	BookingScheduled  = "SCHEDULED"
	BookingInProgress = "IN PROGRESS"
	BookingFinished   = "FINISHED"
)

type Notification struct {
	Type      string `json:"type"`
	Recipient string `json:"recipient"`
}

// BookingSpec defines the desired state of Booking
type BookingSpec struct {
	EndAt         string         `json:"end_at"`
	StartAt       string         `json:"start_at"`
	ResourceName  string         `json:"resource_name"`
	UserID        string         `json:"user_id"`
	Notifications []Notification `json:"notifications,omitempty"`
}

// BookingStatus defines the observed state of Booking
type BookingStatus struct {
	Status           string `json:"status,omitempty"`
	NotificationSent bool   `json:"notification_sent,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:JSONPath=".spec.start_at",name="START",type="string"
//+kubebuilder:printcolumn:JSONPath=".spec.end_at",name="END",type="string"
//+kubebuilder:printcolumn:JSONPath=".status.status",name="STATUS",type="string"

// Booking is the Schema for the bookings API
type Booking struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BookingSpec   `json:"spec,omitempty"`
	Status BookingStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BookingList contains a list of Booking
type BookingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Booking `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Booking{}, &BookingList{})
}
