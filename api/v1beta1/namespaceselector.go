// Copyright Contributors to the Open Cluster Management project

package v1beta1

import (
	"context"
	"fmt"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetNamespaces fetches all namespaces in the cluster and returns a list of the
// namespaces that match the NamespaceSelector. The client.Reader needs access
// for viewing namespaces, like the access given by this kubebuilder tag:
// `//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch`
func (sel NamespaceSelector) GetNamespaces(ctx context.Context, r client.Reader) ([]string, error) {
	matchingNamespaces := make([]string, 0)

	if len(sel.Include) == 0 {
		// A somewhat special case of no matches.
		return matchingNamespaces, nil
	}

	namespaceList := &corev1.NamespaceList{}
	if err := r.List(ctx, namespaceList); err != nil {
		return matchingNamespaces, err
	}

	namespaces := make([]string, len(namespaceList.Items))
	for i, ns := range namespaceList.Items {
		namespaces[i] = ns.GetName()
	}

	return sel.matches(namespaces)
}

// matches filters a slice of strings, and returns ones that match the selector.
// The only possible returned error is a wrapped filepath.ErrBadPattern.
func (sel NamespaceSelector) matches(namespaces []string) ([]string, error) {
	// Using a map to ensure each entry in the result is unique.
	set := make(map[string]struct{})

	for _, namespace := range namespaces {
		include := len(sel.Include) == 0 // include everything if empty/unset

		for _, includePattern := range sel.Include {
			var err error
			include, err = filepath.Match(string(includePattern), namespace)

			if err != nil {
				return []string{}, fmt.Errorf(
					"error parsing 'include' pattern '%s': %w", string(includePattern), err)
			}

			if include {
				break
			}
		}

		if !include {
			continue
		}

		exclude := false

		for _, excludePattern := range sel.Exclude {
			var err error
			exclude, err = filepath.Match(string(excludePattern), namespace)

			if err != nil {
				return []string{}, fmt.Errorf(
					"error parsing 'exclude' pattern '%s': %w", string(excludePattern), err)
			}

			if exclude {
				break
			}
		}

		if exclude {
			continue
		}

		set[namespace] = struct{}{}
	}

	matchingNamespaces := make([]string, 0, len(set))
	for ns := range set {
		matchingNamespaces = append(matchingNamespaces, ns)
	}

	return matchingNamespaces, nil
}
