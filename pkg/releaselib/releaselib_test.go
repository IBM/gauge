package releaselib

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/IBM/gauge/pkg/releaselib"
)

func TestGetLocationMeta(t *testing.T) {
	testlocation := "seattle"
	want := regexp.MustCompile("United States")
	msg, _ := releaselib.GetLocationMeta(testlocation)
	fmt.Println(msg)
	if !want.MatchString(msg) {
		t.Fatalf(`GetLocationMeta("seattle") = %q, want match for %#q, nil`, msg, want)
	}
}

func TestPostLocationMeta(t *testing.T) {
	testlocation := "Flavortown"
	resolved_location := "not resolved"
	err := releaselib.StoreLocationMeta(testlocation, resolved_location)
	if err != nil {
		t.Fatalf(`StoreLocationMeta("Watertown, MA") = %d, not nil`, err)
	}
}
