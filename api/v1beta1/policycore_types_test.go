// Copyright Contributors to the Open Cluster Management project

package v1beta1

import "testing"

func TestIsEnforce(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input RemediationAction
		want  bool
	}{
		{input: "enforce", want: true},
		{input: "Enforce", want: true},
		{input: "ENFORCE", want: false},
		{input: "inform", want: false},
		{input: "Inform", want: false},
		{input: "pleasedoremediate", want: false},
	}

	for _, tc := range tests {
		got := tc.input.IsEnforce()
		if got != tc.want {
			t.Fatalf("Expected IsEnforce to be %v for: '%v', got: %v", tc.want, tc.input, got)
		}
	}
}

func TestIsInform(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input RemediationAction
		want  bool
	}{
		{input: "enforce", want: false},
		{input: "Enforce", want: false},
		{input: "INFORM", want: false},
		{input: "inform", want: true},
		{input: "Inform", want: true},
		{input: "pleasedoremediate", want: false},
	}

	for _, tc := range tests {
		got := tc.input.IsInform()
		if got != tc.want {
			t.Fatalf("Expected IsInform to be %v for: '%v', got: %v", tc.want, tc.input, got)
		}
	}
}
