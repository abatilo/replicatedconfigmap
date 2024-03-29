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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ReplicatedConfigMapSpec defines the desired state of ReplicatedConfigMap
type ReplicatedConfigMapSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Name is the name of the ConfigMap to be created
	Name string `json:"name,omitempty"`

	// Metadata is the metadata to be added to ConfigMaps that are created
	Metadata string `json:"metadata,omitempty"`

	// Data is the data to populate the ConfigMap data key
	Data map[string]string `json:"data,omitempty"`
}

// ReplicatedConfigMapStatus defines the observed state of ReplicatedConfigMap
type ReplicatedConfigMapStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// ReplicatedConfigMap is the Schema for the replicatedconfigmaps API
type ReplicatedConfigMap struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ReplicatedConfigMapSpec   `json:"spec,omitempty"`
	Status ReplicatedConfigMapStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// ReplicatedConfigMapList contains a list of ReplicatedConfigMap
type ReplicatedConfigMapList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ReplicatedConfigMap `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ReplicatedConfigMap{}, &ReplicatedConfigMapList{})
}
