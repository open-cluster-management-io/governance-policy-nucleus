// Copyright Contributors to the Open Cluster Management project

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FakePolicySpec defines the desired state of FakePolicy
type FakePolicySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of FakePolicy. Edit fakepolicy_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// FakePolicyStatus defines the observed state of FakePolicy
type FakePolicyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// FakePolicy is the Schema for the fakepolicies API
type FakePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FakePolicySpec   `json:"spec,omitempty"`
	Status FakePolicyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FakePolicyList contains a list of FakePolicy
type FakePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FakePolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FakePolicy{}, &FakePolicyList{})
}
