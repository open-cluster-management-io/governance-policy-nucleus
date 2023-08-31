// Copyright Contributors to the Open Cluster Management project

package v1beta1

import (
	"context"
	"fmt"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetNamespaces fetches all namespaces in the cluster and returns a list of the
// namespaces that match the NamespaceSelector. The client.Reader needs access
// for viewing namespaces, like the access given by this kubebuilder tag:
// `//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch`
func (sel NamespaceSelector) GetNamespaces(ctx context.Context, r client.Reader) ([]string, error) {
	if len(sel.Include) == 0 && sel.LabelSelector == nil {
		// A somewhat special case of no matches.
		return []string{}, nil
	}

	listOpts := client.ListOptions{}

	if sel.LabelSelector != nil {
		labelSel, err := metav1.LabelSelectorAsSelector(sel.LabelSelector)
		if err != nil {
			return nil, err
		}

		listOpts.LabelSelector = labelSel
	}

	namespaceList := &corev1.NamespaceList{}
	if err := r.List(ctx, namespaceList, &listOpts); err != nil {
		return nil, err
	}

	namespaces := make([]string, len(namespaceList.Items))
	for i, ns := range namespaceList.Items {
		namespaces[i] = ns.GetName()
	}

	return Target(sel).matches(namespaces)
}

type Target struct {
	*metav1.LabelSelector `json:",inline"`

	// Include is a list of filepath expressions to include objects by name.
	Include []NonEmptyString `json:"include,omitempty"`

	// Exclude is a list of filepath expressions to include objects by name.
	Exclude []NonEmptyString `json:"exclude,omitempty"`
}

// matches filters a slice of strings, and returns ones that match the Include
// and Exclude lists in the Target. The only possible returned error is a
// wrapped filepath.ErrBadPattern.
func (t Target) matches(names []string) ([]string, error) {
	// Using a map to ensure each entry in the result is unique.
	set := make(map[string]struct{})

	for _, name := range names {
		include := len(t.Include) == 0 // include everything if empty/unset

		for _, includePattern := range t.Include {
			var err error
			include, err = filepath.Match(string(includePattern), name)

			if err != nil {
				return nil, fmt.Errorf(
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

		for _, excludePattern := range t.Exclude {
			var err error
			exclude, err = filepath.Match(string(excludePattern), name)

			if err != nil {
				return nil, fmt.Errorf(
					"error parsing 'exclude' pattern '%s': %w", string(excludePattern), err)
			}

			if exclude {
				break
			}
		}

		if exclude {
			continue
		}

		set[name] = struct{}{}
	}

	matchingNames := make([]string, 0, len(set))
	for ns := range set {
		matchingNames = append(matchingNames, ns)
	}

	return matchingNames, nil
}
