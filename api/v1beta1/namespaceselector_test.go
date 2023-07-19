// Copyright Contributors to the Open Cluster Management project

package v1beta1

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"k8s.io/apimachinery/pkg/util/validation"
)

var sampleNamespaces = []string{
	"foo", "bar", "baz", "boo", "default", "kube-one", "kube-two", "kube-three",
}

// Fuzz test to verify that excluding "*" always matches 0 namespaces. The
// `Include` list and the input namespaces are both fuzzed.
func FuzzMatchesExcludeAll(f *testing.F) {
	for _, tc := range []string{"", "*", "?", "foo", "quux"} {
		f.Add(tc, tc, tc)
	}

	f.Fuzz(func(t *testing.T, inc1, inc2, extraNS string) {
		errs := validation.IsDNS1123Label(extraNS)
		if len(errs) != 0 { // K8s Namespaces must be valid RFC 1123 DNS labels.
			t.Skip()
		}

		inc := []NonEmptyString{NonEmptyString(inc1), NonEmptyString(inc2)}
		sel := NamespaceSelector{Include: inc, Exclude: []NonEmptyString{"*"}}

		got, err := sel.matches(append(sampleNamespaces, extraNS))
		if err != nil {
			if errors.Is(err, filepath.ErrBadPattern) {
				t.Skip()
			}
			t.Errorf("Unexpected error '%v' when including '%v' and '%v', with ns '%v'",
				err, inc1, inc2, extraNS)
		}
		if len(got) != 0 {
			t.Errorf("Got non-empty matches '%v', when including '%v' and '%v', with ns '%v'",
				got, inc1, inc2, extraNS)
		}
	})
}

// Fuzz test to verify that including "*" will always match all namespaces when
// excludes is empty. Additional items in the `Include` list and the input
// namespaces are both fuzzed.
func FuzzMatchesIncludeAll(f *testing.F) {
	for _, tc := range []string{"", "*", "?", "foo", "quux"} {
		f.Add(tc, tc, tc)
	}

	f.Fuzz(func(t *testing.T, inc1, inc2, extraNS string) {
		errs := validation.IsDNS1123Label(extraNS)
		if len(errs) != 0 { // K8s Namespaces must be valid RFC 1123 DNS labels.
			t.Skip()
		}

		inc := []NonEmptyString{"*", NonEmptyString(inc1), NonEmptyString(inc2)}
		sel := NamespaceSelector{Include: inc, Exclude: []NonEmptyString{}}

		got, err := sel.matches(append(sampleNamespaces, extraNS))
		if err != nil {
			if errors.Is(err, filepath.ErrBadPattern) {
				t.Skip()
			}
			t.Errorf("Unexpected error '%v' when including '%v' and '%v', with ns '%v'",
				err, inc1, inc2, extraNS)
		}

		desiredLen := len(sampleNamespaces) + 1
		for _, ns := range sampleNamespaces {
			if ns == extraNS {
				desiredLen--
			}
		}

		if len(got) != desiredLen {
			t.Errorf("Incorrect matches '%v', when including '%v' and '%v', with ns '%v'",
				got, inc1, inc2, extraNS)
		}
	})
}

func TestMatches(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		inc  []NonEmptyString
		exc  []NonEmptyString
		want []string
	}{
		"include trailing wildcard": {
			inc:  []NonEmptyString{"kube-*"},
			exc:  []NonEmptyString{},
			want: []string{"kube-one", "kube-two", "kube-three"},
		},
		"exclude trailing wildcard": {
			inc:  []NonEmptyString{"*"},
			exc:  []NonEmptyString{"kube-*"},
			want: []string{"foo", "bar", "baz", "boo", "default"},
		},
		"include leading wildcard": {
			inc:  []NonEmptyString{"*o"},
			exc:  []NonEmptyString{},
			want: []string{"foo", "boo", "kube-two"},
		},
		"exclude leading wildcard": {
			inc:  []NonEmptyString{"*"},
			exc:  []NonEmptyString{"*o"},
			want: []string{"bar", "baz", "default", "kube-one", "kube-three"},
		},
		"include middle wildcard": {
			inc:  []NonEmptyString{"*a*"},
			exc:  []NonEmptyString{},
			want: []string{"bar", "baz", "default"},
		},
		"exclude middle wildcard": {
			inc:  []NonEmptyString{"*"},
			exc:  []NonEmptyString{"*a*"},
			want: []string{"foo", "boo", "kube-one", "kube-two", "kube-three"},
		},
		"include one specific": {
			inc:  []NonEmptyString{"foo"},
			exc:  []NonEmptyString{},
			want: []string{"foo"},
		},
		"exclude one specific": {
			inc:  []NonEmptyString{"*"},
			exc:  []NonEmptyString{"foo"},
			want: []string{"bar", "baz", "boo", "default", "kube-one", "kube-two", "kube-three"},
		},
		"include multiple specific": {
			inc:  []NonEmptyString{"foo", "bar", "kube-three"},
			exc:  []NonEmptyString{},
			want: []string{"foo", "bar", "kube-three"},
		},
		"exclude multiple specific": {
			inc:  []NonEmptyString{"*"},
			exc:  []NonEmptyString{"foo", "bar", "kube-three"},
			want: []string{"baz", "boo", "default", "kube-one", "kube-two"},
		},
		"include single char wildcards": {
			inc:  []NonEmptyString{"?a?"},
			exc:  []NonEmptyString{},
			want: []string{"bar", "baz"},
		},
		"exclude single char wildcards": {
			inc:  []NonEmptyString{"*"},
			exc:  []NonEmptyString{"?a?"},
			want: []string{"foo", "boo", "default", "kube-one", "kube-two", "kube-three"},
		},
		"include character class": {
			inc:  []NonEmptyString{"[fb]oo"},
			exc:  []NonEmptyString{},
			want: []string{"foo", "boo"},
		},
		"exclude character class": {
			inc:  []NonEmptyString{"*"},
			exc:  []NonEmptyString{"[fb]oo"},
			want: []string{"bar", "baz", "default", "kube-one", "kube-two", "kube-three"},
		},
	}

	for name, tcase := range tests {
		sel := NamespaceSelector{Include: tcase.inc, Exclude: tcase.exc}

		got, err := sel.matches(sampleNamespaces)
		if err != nil {
			t.Errorf("Unexpected error '%v', in test '%v'", err, name)
		}

		less := func(a, b string) bool { return a < b }
		diff := cmp.Diff(tcase.want, got, cmpopts.SortSlices(less))

		if diff != "" {
			t.Errorf("Mismatch in test '%v': %v", name, diff)
		}
	}
}
