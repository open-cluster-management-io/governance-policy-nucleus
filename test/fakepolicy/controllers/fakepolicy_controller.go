// Copyright Contributors to the Open Cluster Management project

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	policyv1beta1 "open-cluster-management.io/governance-policy-nucleus/test/fakepolicy/api/v1beta1"
)

// FakePolicyReconciler reconciles a FakePolicy object
type FakePolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Usual RBAC for fakepolicy:
//+kubebuilder:rbac:groups=policy.open-cluster-management.io,resources=fakepolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=policy.open-cluster-management.io,resources=fakepolicies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=policy.open-cluster-management.io,resources=fakepolicies/finalizers,verbs=update

// Nucleus RBAC for namespaceSelector:
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch

func (r *FakePolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	policy := &policyv1beta1.FakePolicy{}
	if err := r.Get(ctx, req.NamespacedName, policy); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, probably deleted
			return ctrl.Result{}, nil
		}

		log.Error(err, "Failed to get FakePolicy")

		return ctrl.Result{}, err
	}

	selectedNamespaces, err := policy.Spec.NamespaceSelector.GetNamespaces(ctx, r.Client)
	if err != nil {
		log.Error(err, "Failed to GetNamespaces using NamespaceSelector",
			"selector", policy.Spec.NamespaceSelector)

		policy.Status.SelectionError = err.Error()
	} else {
		policy.Status.SelectionError = ""
	}

	policy.Status.SelectedNamespaces = selectedNamespaces
	policy.Status.SelectionComplete = true

	if err := r.Status().Update(ctx, policy); err != nil {
		log.Error(err, "Failed to update status")

		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FakePolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&policyv1beta1.FakePolicy{}).
		Complete(r)
}
