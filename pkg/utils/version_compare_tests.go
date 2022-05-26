package utils

import (
	"testing"
)

func TestEquals(t *testing.T) {
	res := IsEqual("1.0.2", "1.0.2")
	if !res {
		t.Fail()
	}
	res = IsEqual("1.0.2", "1.0.3")
	if res {
		t.Fail()
	}
	res = IsEqual("v1.0.2", "v1.0.2")
	if res {
		t.Fail()
	}
}
