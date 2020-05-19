/*


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

// TrackingServerSpec defines the desired state of TrackingServer
// +k8s:openapi-gen=true
type TrackingServerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Size          int32  `json:"size"`
	Image         string `json:"image"`
	S3endpointURL string `json:"S3endpointURL,omitempty"`
	S3secretName  string `json:"S3secretName,omitempty"`
}

// TrackingServerStatus defines the observed state of TrackingServer
// +k8s:openapi-gen=true
type TrackingServerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TrackingServer is the Schema for the trackingservers API
// +k8s:openapi-gen=true
type TrackingServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TrackingServerSpec   `json:"spec,omitempty"`
	Status TrackingServerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TrackingServerList contains a list of TrackingServer
type TrackingServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TrackingServer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TrackingServer{}, &TrackingServerList{})
}
