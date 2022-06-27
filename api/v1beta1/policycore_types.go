// Copyright Contributors to the Open Cluster Management project

// Package v1beta1 contains API Schema definitions for the policy v1beta1 API group
//+kubebuilder:object:generate=true
//+groupName=policy.open-cluster-management.io
package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// Important: Run "make" to regenerate code after modifying this file

// PolicyCoreSpec defines the desired state of PolicyCore
type PolicyCoreSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster

	// Foo is an example field of PolicyCore. Edit policycore_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// PolicyCoreStatus defines the observed state of PolicyCore
type PolicyCoreStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
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
