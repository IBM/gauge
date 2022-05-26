package utils

import (
	"github.com/Masterminds/semver"
)

//IsEqual :
func IsEqual(base, target string) bool {
	baseV, err := semver.NewVersion(base)
	if err != nil {
		// fmt.Printf("error parsing version `%v`: %v\n", base, err)
		return false
	}
	targetV, err := semver.NewVersion(target)
	if err != nil {
		// fmt.Printf("error parsing version `%v`: %v\n", target, err)
		return false
	}
	return baseV.Equal(targetV)
}

func IsGreater(base, target string) bool {
	baseV, err := semver.NewVersion(base)
	if err != nil {
		// fmt.Printf("error parsing version `%v`: %v\n", base, err)
		return false
	}
	targetV, err := semver.NewVersion(target)
	if err != nil {
		// fmt.Printf("error parsing version `%v`: %v\n", target, err)
		return false
	}
	return baseV.GreaterThan(targetV)
}

func IsGreaterMajor(base, target string) bool {
	baseV, err := semver.NewVersion(base)
	if err != nil {
		// fmt.Printf("error parsing version `%v`: %v\n", base, err)
		return false
	}
	targetV, err := semver.NewVersion(target)
	if err != nil {
		// fmt.Printf("error parsing version `%v`: %v\n", target, err)
		return false
	}
	return baseV.Major() > targetV.Major()
}
