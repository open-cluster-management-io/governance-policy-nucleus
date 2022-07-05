// Copyright Contributors to the Open Cluster Management project

// Package v1beta1 contains API Schema definitions for the policy v1beta1 API group
//+kubebuilder:object:generate=true
//+groupName=policy.open-cluster-management.io
//+kubebuilder:validation:Optional
package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required. Any new fields you add must have json tags for
// the fields to be serialized.
// Important: Run "make" to regenerate code after modifying this file

// PolicyCoreSpec defines fields that policies must implement to be part of the
// Open Cluster Management policy framework. The intention is for controllers
// to embed this struct in their *Spec definitions.
type PolicyCoreSpec struct {
	// Severity defines how serious the situation is when the policy is not
	// compliant. The severity should not change the behavior of the policy, but
	// may be read and used by other tools. Accepted values include: low,
	// medium, high, and critical.
	Severity Severity `json:"severity,omitempty"`

	// RemediationAction indicates what the policy controller should do when the
	// policy is not compliant. Accepted values include inform, and enforce.
	// Note that not all policy controllers will attempt to automatically
	// remediate a policy, even when set to "enforce".
	RemediationAction RemediationAction `json:"remediationAction,omitempty"`

	// NamepaceSelector indicates which namespaces on the cluster this policy
	// should apply to, when the policy applies to namespaced objects.
	NamespaceSelector NamespaceSelector `json:"namespaceSelector,omitempty"`
}

//+kubebuilder:validation:Enum=low;Low;medium;Medium;high;High;critical;Critical
type Severity string

//+kubebuilder:validation:Enum=Inform;inform;Enforce;enforce
type RemediationAction string

// IsEnforce is true when the policy controller can attempt to enforce the
// policy by remediating it automatically. Note that not all controllers will
// support automatic enforcement.
func (ra RemediationAction) IsEnforce() bool {
	return ra == "Enforce" || ra == "enforce"
}

// IsInform is true when the policy controller should only report whether the
// policy is compliant or not and should not perform any actions to attempt
// remediation.
func (ra RemediationAction) IsInform() bool {
	return ra == "Inform" || ra == "inform"
}

type NamespaceSelector struct {
	// FUTURE: enhance this per
	// https://github.com/open-cluster-management-io/enhancements/tree/main/enhancements/sig-policy/62-namespace-labelselector

	// Include is a list of namespaces the policy should apply to.
	Include []NonEmptyString `json:"include,omitempty"`

	// Exclude is a list of namespaces the policy should _not_ apply to.
	Exclude []NonEmptyString `json:"exclude,omitempty"`
}

//+kubebuilder:validation:MinLength=1
type NonEmptyString string

// PolicyCoreStatus defines fields that policies should implement as part of
// the Open Cluster Management policy framework.
type PolicyCoreStatus struct {
	// ComplianceState indicates whether the policy is compliant or not.
	// Accepted values include: Compliant, NonCompliant, and UnknownCompliancy
	ComplianceState ComplianceState `json:"compliant,omitempty"`
}

//+kubebuilder:validation:Enum=Compliant;NonCompliant;UnknownCompliancy
type ComplianceState string

const (
	// Compliant indicates that the policy controller determined there were no
	// violations to the policy in the cluster.
	Compliant ComplianceState = "Compliant"

	// NonCompliant indicates that the policy controller found an issue in the
	// cluster that is considered a violation.
	NonCompliant ComplianceState = "NonCompliant"

	// UnknownCompliancy indicates that the policy controller could not determine
	// if the cluster has any violations or not.
	UnknownCompliancy ComplianceState = "UnknownCompliancy"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PolicyCore is the Schema for the policycores API
type PolicyCore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PolicyCoreSpec   `json:"spec,omitempty"`
	Status PolicyCoreStatus `json:"status,omitempty"`
}
