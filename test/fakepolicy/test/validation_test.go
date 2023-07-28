// Copyright Contributors to the Open Cluster Management project

package test

import (
	"context"
	"embed"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"

	policyv1beta1 "open-cluster-management.io/governance-policy-nucleus/test/fakepolicy/api/v1beta1"
)

var _ = Describe("FakePolicy CRD Validation", Ordered, func() {
	AfterEach(func() {
		By("Deleting any FakePolicies in the default namespace")
		fakePolicy := &policyv1beta1.FakePolicy{}
		Expect(k8sClient.DeleteAllOf(context.TODO(), fakePolicy, &client.DeleteAllOfOptions{
			ListOptions: client.ListOptions{Namespace: "default"},
		})).Should(Succeed())
	})

	DescribeTable("Validating spec inputs",
		func(severity, remediationAction string, include, exclude []string, isValid bool) {
			policy, nn := fromTestdata("fakepolicy-sample.yaml")

			Expect(unstructured.SetNestedField(policy.Object,
				severity, "spec", "severity")).To(Succeed())
			Expect(unstructured.SetNestedField(policy.Object,
				remediationAction, "spec", "remediationAction")).To(Succeed())
			Expect(unstructured.SetNestedStringSlice(policy.Object,
				include, "spec", "namespaceSelector", "include")).To(Succeed())
			Expect(unstructured.SetNestedStringSlice(policy.Object,
				exclude, "spec", "namespaceSelector", "exclude")).To(Succeed())

			matchExpected := Succeed
			if !isValid {
				matchExpected = HaveOccurred
			}

			Expect(k8sClient.Create(context.TODO(), policy)).Should(matchExpected())

			foundPolicy := &policyv1beta1.FakePolicy{}
			Expect(k8sClient.Get(context.TODO(), nn, foundPolicy)).Should(matchExpected())
		},
		// Test severity options
		Entry("severity=low", "low", "inform", []string{"*"}, []string{"kube-*"}, true),
		Entry("severity=medium", "medium", "inform", []string{"*"}, []string{"kube-*"}, true),
		Entry("severity=high", "high", "inform", []string{"*"}, []string{"kube-*"}, true),
		Entry("severity=critical", "critical", "inform", []string{"*"}, []string{"kube-*"}, true),
		Entry("severity=Low", "Low", "inform", []string{"*"}, []string{"kube-*"}, true),
		Entry("severity=LOW", "LOW", "inform", []string{"*"}, []string{"kube-*"}, false),
		Entry("severity=super", "super", "inform", []string{"*"}, []string{"kube-*"}, false),
		Entry("severity=''", "", "inform", []string{"*"}, []string{"kube-*"}, false),

		// Test remediationAction options
		Entry("remediationAction=inform", "low", "inform", []string{"*"}, []string{"kube-*"}, true),
		Entry("remediationAction=enforce", "low", "enforce", []string{"*"}, []string{"kube-*"}, true),
		Entry("remediationAction=Inform", "low", "Inform", []string{"*"}, []string{"kube-*"}, true),
		Entry("remediationAction=Enforce", "low", "Enforce", []string{"*"}, []string{"kube-*"}, true),
		Entry("remediationAction=INFORM", "low", "INFORM", []string{"*"}, []string{"kube-*"}, false),
		Entry("remediationAction=like-magic", "low", "like-magic", []string{"*"}, []string{"kube-*"}, false),

		// Test namespaceSelector options
		Entry("empty string in namespaceSelector.include", "low", "inform", []string{""}, []string{"kube-*"}, false),
		Entry("empty string in namespaceSelector.exclude", "low", "inform", []string{"*"}, []string{""}, false),
		Entry("empty list in namespaceSelector.include", "low", "inform", []string{}, []string{"kube-*"}, true),
		Entry("empty list in namespaceSelector.exclude", "low", "inform", []string{"*"}, []string{}, true),
	)
})

//go:embed testdata/*
var testfiles embed.FS

// Unmarshals the given YAML file in testdata/ into an unstructured.Unstructured,
// and additionally returns a NamespacedName for easier lookup later.
func fromTestdata(name string) (*unstructured.Unstructured, types.NamespacedName) {
	objYAML, err := testfiles.ReadFile("testdata/" + name)
	Expect(err).ToNot(HaveOccurred())

	m := make(map[string]interface{})
	Expect(yaml.UnmarshalStrict(objYAML, &m)).To(Succeed())

	unstruct := &unstructured.Unstructured{Object: m}
	nn := types.NamespacedName{
		Namespace: unstruct.GetNamespace(),
		Name:      unstruct.GetName(),
	}

	return unstruct, nn
}
