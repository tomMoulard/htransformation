package assert

import (
	"fmt"
	"reflect"
	"testing"
)

func Error(t *testing.T, err error) bool {
	t.Helper()

	if err == nil {
		t.Logf("Error is nil, but should have an error")
		t.Fail()

		return false
	}

	return true
}

func NoError(t *testing.T, err error) bool {
	t.Helper()

	if err != nil {
		t.Logf("Error is not nil, but should not have an error")
		t.Fail()

		return false
	}

	return true
}

func Equal(t *testing.T, expect, actual interface{}) bool {
	t.Helper()

	if !reflect.DeepEqual(expect, actual) {
		t.Logf("Expect %v, but got %v", expect, actual)
		t.Fail()

		return false
	}

	return true
}

func Equalf(t *testing.T, expect, actual interface{}, format string, args ...interface{}) bool {
	t.Helper()

	if !reflect.DeepEqual(expect, actual) {
		t.Logf("Expect %v, but got %v: %s", expect, actual, fmt.Sprintf(format, args...))
		t.Fail()

		return false
	}

	return true
}
