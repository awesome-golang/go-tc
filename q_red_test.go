package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRed(t *testing.T) {
	tests := map[string]struct {
		val  Red
		err1 error
		err2 error
	}{
		"simple": {val: Red{MaxP: 2, Parms: &RedQOpt{QthMin: 2, QthMax: 4}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalRed(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Red{}
			err2 := unmarshalRed(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Red missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalRed(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
