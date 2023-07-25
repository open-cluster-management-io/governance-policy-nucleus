// Copyright Contributors to the Open Cluster Management project

package controllers

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	nucleusv1beta1 "open-cluster-management.io/governance-policy-nucleus/api/v1beta1"
	policyv1beta1 "open-cluster-management.io/governance-policy-nucleus/test/fakepolicy/api/v1beta1"
)

// Most tests of the NamespaceSelector should be done in unit tests. These are
// just some additional sanity checks.
var _ = Describe("FakePolicy NamespaceSelection", Ordered, func() {
	defaultNamespaces := []string{"default", "kube-node-lease", "kube-public", "kube-system"}
	sampleNamespaces := []string{"foo", "goo", "fake", "faze", "kube-one"}
	allNamespaces := append(defaultNamespaces, sampleNamespaces...)

	BeforeAll(func() {
		By("Creating sample namespaces")
		for _, ns := range sampleNamespaces {
			nsObj := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}}
			Expect(k8sClient.Create(context.TODO(), nsObj)).To(Succeed())
		}

		By("Ensuring the allNamespaces list is correct")
		// constructing the default / allNamespaces lists is complicated because of how ginkgo
		// runs the table tests... this seems better than other workarounds.
		nsList := corev1.NamespaceList{}
		Expect(k8sClient.List(context.TODO(), &nsList)).To(Succeed())

		foundNS := make([]string, len(nsList.Items))
		for i, ns := range nsList.Items {
			foundNS[i] = ns.GetName()
		}

		Expect(allNamespaces).To(ConsistOf(foundNS))
	})

	AfterAll(func() {
		By("Deleting sample namespaces")
		for _, ns := range sampleNamespaces {
			nsObj := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}}
			Expect(k8sClient.Delete(context.TODO(), nsObj)).To(Succeed())
		}
	})

	AfterEach(func() {
		By("Deleting any FakePolicies in the default namespace")
		fakePolicy := &policyv1beta1.FakePolicy{}
		Expect(k8sClient.DeleteAllOf(context.TODO(), fakePolicy, &client.DeleteAllOfOptions{
			ListOptions: client.ListOptions{Namespace: "default"},
		})).To(Succeed())
	})

	DescribeTable("Verifying NamespaceSelector behavior",
		func(sel nucleusv1beta1.NamespaceSelector, desiredMatches []string, selErr string) {
			policy := sampleFakePolicy()
			policy.Spec.NamespaceSelector = sel
			Expect(k8sClient.Create(context.TODO(), &policy)).To(Succeed())

			nn := types.NamespacedName{
				Name:      policy.GetName(),
				Namespace: policy.GetNamespace(),
			}

			Eventually(func(g Gomega) {
				foundPolicy := policyv1beta1.FakePolicy{}
				g.Expect(k8sClient.Get(context.TODO(), nn, &foundPolicy)).To(Succeed())

				g.Expect(foundPolicy.Status.SelectionComplete).To(BeTrue())
				g.Expect(foundPolicy.Status.SelectedNamespaces).To(ConsistOf(desiredMatches))
				g.Expect(foundPolicy.Status.SelectionError).To(Equal(selErr))
			}).Should(Succeed())
		},

		// Basic testing of includes and excludes
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

func sampleFakePolicy() policyv1beta1.FakePolicy {
	return policyv1beta1.FakePolicy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "policy.open-cluster-management.io/v1beta1",
			Kind:       "FakePolicy",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fakepolicy-sample",
			Namespace: "default",
		},
		Spec: policyv1beta1.FakePolicySpec{
			PolicyCoreSpec: nucleusv1beta1.PolicyCoreSpec{
				Severity:          "low",
				RemediationAction: "inform",
				NamespaceSelector: nucleusv1beta1.NamespaceSelector{
					Include: []nucleusv1beta1.NonEmptyString{"*"},
					Exclude: []nucleusv1beta1.NonEmptyString{"kube-*"},
				},
			},
		},
	}
}
