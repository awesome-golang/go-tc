package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestHhf(t *testing.T) {
	tests := map[string]struct {
		val  Hhf
		err1 error
		err2 error
	}{
		"simple": {val: Hhf{BacklogLimit: 1, Quantum: 2, HHFlowsLimit: 3, ResetTimeout: 4, AdmitBytes: 5, EVICTTimeout: 6, NonHHWeight: 7}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalHhf(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Hhf{}
			err2 := unmarshalHhf(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Hhf missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalHhf(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
