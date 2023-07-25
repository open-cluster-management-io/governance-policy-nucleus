//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Copyright Contributors to the Open Cluster Management project

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FakePolicy) DeepCopyInto(out *FakePolicy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FakePolicy.
func (in *FakePolicy) DeepCopy() *FakePolicy {
	if in == nil {
		return nil
	}
	out := new(FakePolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FakePolicy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FakePolicyList) DeepCopyInto(out *FakePolicyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FakePolicy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FakePolicyList.
func (in *FakePolicyList) DeepCopy() *FakePolicyList {
	if in == nil {
		return nil
	}
	out := new(FakePolicyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FakePolicyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FakePolicySpec) DeepCopyInto(out *FakePolicySpec) {
	*out = *in
	in.PolicyCoreSpec.DeepCopyInto(&out.PolicyCoreSpec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FakePolicySpec.
func (in *FakePolicySpec) DeepCopy() *FakePolicySpec {
	if in == nil {
		return nil
	}
	out := new(FakePolicySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FakePolicyStatus) DeepCopyInto(out *FakePolicyStatus) {
	*out = *in
	out.PolicyCoreStatus = in.PolicyCoreStatus
	if in.SelectedNamespaces != nil {
		in, out := &in.SelectedNamespaces, &out.SelectedNamespaces
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FakePolicyStatus.
func (in *FakePolicyStatus) DeepCopy() *FakePolicyStatus {
	if in == nil {
		return nil
	}
	out := new(FakePolicyStatus)
	in.DeepCopyInto(out)
	return out
}
