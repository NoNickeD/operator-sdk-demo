/*
Copyright 2025.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PodNotifRestartSpec defines the desired state of PodNotifRestart
type PodNotifRestartSpec struct {
	MinRestarts int32 `json:"minRestarts"`
}

// PodNotifRestartStatus defines the observed state of PodNotifRestart
type PodNotifRestartStatus struct {
	LastNotification string `json:"lastNotification,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// PodNotifRestart is the Schema for the podnotifrestarts API
type PodNotifRestart struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodNotifRestartSpec   `json:"spec,omitempty"`
	Status PodNotifRestartStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PodNotifRestartList contains a list of PodNotifRestart
type PodNotifRestartList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodNotifRestart `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PodNotifRestart{}, &PodNotifRestartList{})
}
