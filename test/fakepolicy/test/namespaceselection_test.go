// Copyright Contributors to the Open Cluster Management project

package test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	nucleusv1beta1 "open-cluster-management.io/governance-policy-nucleus/api/v1beta1"
	policyv1beta1 "open-cluster-management.io/governance-policy-nucleus/test/fakepolicy/api/v1beta1"
)

var _ = Describe("FakePolicy NamespaceSelection", Ordered, func() {
	defaultNamespaces := []string{"default", "kube-node-lease", "kube-public", "kube-system"}
	sampleNamespaces := []string{"foo", "goo", "fake", "faze", "kube-one"}
	allNamespaces := append(defaultNamespaces, sampleNamespaces...)

	BeforeAll(func() {
		By("Creating sample namespaces")
		for _, ns := range sampleNamespaces {
			nsObj := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}}
			Expect(cleanlyCreate(nsObj)).To(Succeed())
		}

		By("Ensuring the allNamespaces list is correct")
		// constructing the default / allNamespaces lists is complicated because of how ginkgo
		// runs the table tests... this seems better than other workarounds.
		nsList := corev1.NamespaceList{}
		Expect(k8sClient.List(ctx, &nsList)).To(Succeed())

		foundNS := make([]string, len(nsList.Items))
		for i, ns := range nsList.Items {
			foundNS[i] = ns.GetName()
		}

		Expect(allNamespaces).To(ConsistOf(foundNS))
	})

	DescribeTable("Verifying NamespaceSelector behavior",
		func(sel nucleusv1beta1.NamespaceSelector, desiredMatches []string, selErr string) {
			policy := sampleFakePolicy()
			policy.Spec.NamespaceSelector = sel

			Expect(cleanlyCreate(&policy)).To(Succeed())

			Eventually(func(g Gomega) {
				foundPolicy := policyv1beta1.FakePolicy{}
				g.Expect(k8sClient.Get(ctx, getNamespacedName(&policy), &foundPolicy)).To(Succeed())
				g.Expect(foundPolicy.Status.SelectionComplete).To(BeTrue())
				g.Expect(foundPolicy.Status.SelectedNamespaces).To(ConsistOf(desiredMatches))
				g.Expect(foundPolicy.Status.SelectionError).To(Equal(selErr))
			}).Should(Succeed())
		},

		// Basic testing of includes and excludes
		Entry("empty should match no namespaces", nucleusv1beta1.NamespaceSelector{},
			[]string{}, ""),
		Entry("include all with *", nucleusv1beta1.NamespaceSelector{
			Include: []nucleusv1beta1.NonEmptyString{"*"},
		}, allNamespaces, ""),
		Entry("include a specific namespace", nucleusv1beta1.NamespaceSelector{
			Include: []nucleusv1beta1.NonEmptyString{"foo"},
		}, []string{"foo"}, ""),
		Entry("include multiple namespaces with a wildcard", nucleusv1beta1.NamespaceSelector{
			Include: []nucleusv1beta1.NonEmptyString{"fa?e"},
		}, []string{"fake", "faze"}, ""),
		Entry("exclude namespaces with a wildcard", nucleusv1beta1.NamespaceSelector{
			Include: []nucleusv1beta1.NonEmptyString{"*"},
			Exclude: []nucleusv1beta1.NonEmptyString{"kube-*"},
		}, []string{"default", "foo", "goo", "fake", "faze"}, ""),
		Entry("error if an include entry is malformed", nucleusv1beta1.NamespaceSelector{
			Include: []nucleusv1beta1.NonEmptyString{"kube-[system"},
		}, []string{}, "error parsing 'include' pattern 'kube-[system': syntax error in pattern"),
	)
})
