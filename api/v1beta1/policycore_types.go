// Copyright Contributors to the Open Cluster Management project

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PolicyCoreSpec defines the desired state of PolicyCore
type PolicyCoreSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of PolicyCore. Edit policycore_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// PolicyCoreStatus defines the observed state of PolicyCore
type PolicyCoreStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PolicyCore is the Schema for the policycores API
type PolicyCore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PolicyCoreSpec   `json:"spec,omitempty"`
	Status PolicyCoreStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PolicyCoreList contains a list of PolicyCore
type PolicyCoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PolicyCore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PolicyCore{}, &PolicyCoreList{})
}
