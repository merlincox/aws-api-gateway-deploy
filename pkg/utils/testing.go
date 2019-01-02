package utils

import "testing"

func AssertEquals(t *testing.T, context string, expected, result interface{}) {
	if expected != result {
		t.Fatalf("%v should be '%v', but was '%v'", context, expected, result)
	}
}

func AssertErrorEquals(t *testing.T, context string, expected, result interface{}) {

	if result == nil {
		t.Fatalf("%v should be '%v', but was nil", context, expected)
	}

	err, ok := result.(error)

	if !ok {
		t.Fatalf("%v should be error '%v', but was non-error '%v'", context, expected, result)
	}

	if expected != err.Error() {
		t.Fatalf("%v should be '%v', but was '%v'", context, expected, err.Error())
	}
}

func AssertNoError(t *testing.T, context string, err error) {
	if err != nil {
		t.Fatalf("%v should throw no error, but threw '%v'", context, err)
	}
}

func AssertTrue(t *testing.T, context string, test bool) {
	if test != true {
		t.Fatalf("%v should be true but was false", context)
	}
}

func AssertFalse(t *testing.T, context string, test bool) {
	if test == true {
		t.Fatalf("%v should be false but was true", context)
	}
}
