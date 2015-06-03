package ae

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestErrors(t *testing.T) {
	simpleErr := errors.New("test")
	fmtErr := fmt.Errorf("err %v %v %v", 1, 2, 4)

	tests := []struct {
		name      string
		err       Err
		wantErrs  []error
		msgPrefix string
	}{
		{
			name:      "Wrap simple error",
			err:       Wrap(simpleErr),
			wantErrs:  []error{simpleErr},
			msgPrefix: simpleErr.Error(),
		},
		{
			name:      "Wrapf simple error",
			err:       Wrapf(simpleErr, "err %v %v %v", 1, 2, 4),
			wantErrs:  []error{simpleErr, fmtErr},
			msgPrefix: fmtErr.Error() + ": " + simpleErr.Error(),
		},
		{
			name:      "Errorf",
			err:       Errorf("err %v %v %v", 1, 2, 4),
			wantErrs:  []error{fmtErr},
			msgPrefix: fmtErr.Error(),
		},
		{
			name:      "Wrap ae.Err",
			err:       Wrap(Wrap(simpleErr)),
			wantErrs:  []error{simpleErr},
			msgPrefix: simpleErr.Error(),
		},
		{
			name:      "Wrapf ae.Err",
			err:       Wrapf(Wrap(simpleErr), "err %v %v %v", 1, 2, 4),
			wantErrs:  []error{simpleErr, fmtErr},
			msgPrefix: fmtErr.Error() + ": " + simpleErr.Error(),
		},
	}

	for _, tt := range tests {
		if !reflect.DeepEqual(tt.err.Errors(), tt.wantErrs) {
			t.Errorf("%s: Errors() failed:\n  got %v\n want %v", tt.name, tt.err.Errors(), tt.wantErrs)
		}

		if !strings.HasPrefix(tt.err.Error(), tt.msgPrefix) {
			t.Errorf("%s: Error() failed:\n  got: %v\n want prefix %v", tt.name, tt.err.Error(), tt.msgPrefix)
		}

		stackParts := strings.Split(tt.err.Stack(), StackSeparator)
		if len(stackParts) == 0 {
			t.Fatalf("Stack is missing: %v", tt.err.Stack())
		}

		if !strings.Contains(stackParts[0], "ae.TestErrors: ") {
			t.Errorf("s[0] should be ae.TestErrors, got %v", stackParts[0])
		}
	}
}

func TestWrapNested(t *testing.T) {
	const testCalls = 200
	var nested func(calls int) Err
	nested = func(calls int) Err {
		if calls == 0 {
			return Errorf("nested")
		}
		return nested(calls - 1)
	}

	got := nested(testCalls)
	stackParts := strings.Split(got.Stack(), StackSeparator)
	if len(stackParts) < testCalls {
		t.Errorf("Stack does not contain all nested calls, got %v, expected >%v", len(stackParts), testCalls)
	}
}

func TestWrapNil(t *testing.T) {
	if got := Wrap(nil); got != nil {
		t.Errorf("Wrap(nil) should return nil, got %v", got)
	}
	if got := Wrapf(nil, "message %v", 1); got != nil {
		t.Errorf("Wrapf(nil, ...) should return nil, got %v", got)
	}
}
