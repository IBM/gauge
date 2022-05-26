package core

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/IBM/gauge/pkg/common"
	"gopkg.in/yaml.v2"
)

type gaugeControl struct {
	RuntimeConfig struct {
		ReleaseLibService string `yaml:"releaselib-service"`
	} `yaml:"runtime-config"`

	ReleaseConfig struct {
		Enable                bool `yaml:"enable"`
		MaxReleaseLag         int  `yaml:"max-release-lag"`
		MaxReleaseLagDuration int  `yaml:"max-release-lag-duration"`
		PeerReviewEnforced    bool `yaml:"peer-review-enforced"`
		ZombieCommitEnforced  bool `yaml:"zombie-commit-enforced"`
	} `yaml:"release-control"`

	ExportControl struct {
		Enable                bool     `yaml:"enable"`
		RestrictedCountries   []string `yaml:"restricted-countries"`
		TAACountries          []string `yaml:"taa-list"`
		OFACCountries         []string `yaml:"ofac-list"`
		ContributionThreshold int      `yaml:"contribution-threshold"`
	} `yaml:"export-controls"`
}

func parseGaugeControls(filepath string, controlOpts *gaugeControl) error {
	yamlBuf, err := ioutil.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("error reading control file: %s: %v", filepath, err)
	}
	if err := yaml.Unmarshal(yamlBuf, &controlOpts); err != nil {
		return fmt.Errorf("error parsing control file: %s: %v", filepath, err)
	}
	return nil
}

// func evaluateOSSDependency(packageHealth common.PackageHealthStat, controlOpts *gaugeControl) (common.FailePackageReport, common.ErrorPackageReport, bool) {
// 	failReport := common.FailePackageReport{}
// 	failReport.Name = packageHealth.Name
// 	failReport.Version = packageHealth.Version
// 	checkOK := true
// 	failReport.TotalChecks = len(packageHealth.Contributors)

// 	errorReport := common.ErrorPackageReport{}
// 	errorReport.Name = packageHealth.Name
// 	errorReport.Version = packageHealth.Version
// 	errorReport.TotalChecks = len(packageHealth.Contributors)

// 	//1. check contributor's location for ExportLocationRestrictions

// 	// following map/find approach did not work, because the
// 	// country name might not be exact match. FIXME

// 	// rlocs := map[string]struct{}{}
// 	// for _, loc := range controlOpts.LicenseRestrictions {
// 	// 	rlocs[strings.ToLower(loc)] = struct{}{}
// 	// }
// 	for _, cloc := range packageHealth.Contributors {
// 		cleanLocation := weather.RemoveCharacters(cloc.Location, "<>;&|\"{}\\^%")

// 		if cleanLocation == "not available" || cleanLocation == "navail" || cleanLocation == "" {
// 			// fmt.Printf("no location data in profile for contributor: \"%s\", URL: \"%s\"\n", cloc.Name, cloc.URL)
// 			errorReport.Reasons = append(errorReport.Reasons, fmt.Sprintf("no location data in profile for contributor: \"%s\", URL: \"%s\"", cloc.Name, cloc.URL))
// 			errorReport.ErrorChecks++
// 			continue
// 		}
// 		resolvedCountry, err := releaselib.GetLocationMeta(cleanLocation)
// 		if err != nil {
// 			resolvedCountry, err = weather.GetLocation(cleanLocation)
// 			if err != nil {
// 				log.Fatalf("error resolving location \"%s\": %s ", cloc.Location, err)
// 			}

// 			if resolvedCountry == "" {
// 				resolvedCountry = "not resolved"
// 			}
// 			// releaselib.StoreLocationMeta(cleanLocation, resolvedCountry)
// 		}
// 		// fmt.Printf("TEST: location: %s, resolved location: %s\n", cloc.Location, resolvedCountry)
// 		if resolvedCountry == "not resolved" {
// 			// fmt.Printf("location data not resolved for contributor: \"%s\", URL: \"%s\", location: \"%s\"\n", cloc.Name, cloc.URL, cloc.Location)
// 			errorReport.Reasons = append(errorReport.Reasons, fmt.Sprintf("location data not resolved for contributor: \"%s\", URL: \"%s\", location: \"%s\"", cloc.Name, cloc.URL, cloc.Location))
// 			errorReport.ErrorChecks++
// 			continue
// 		}

// 		for _, rloc := range controlOpts.ExportControl.RestrictedCountries {
// 			if strings.Contains(strings.ToLower(resolvedCountry), strings.ToLower(rloc)) {
// 				// fmt.Printf("export_control_location check failed for contributor: \"%s\", URL: \"%s\", location: \"%s\", resolved country: \"%s\"\n", cloc.Name, cloc.URL, cloc.Location, resolvedCountry)
// 				failReport.Reasons = append(failReport.Reasons, fmt.Sprintf("export_control_location check failed for contributor: \"%s\", URL: \"%s\", location: \"%s\", resolved country: \"%s\"", cloc.Name, cloc.URL, cloc.Location, resolvedCountry))
// 				checkOK = false
// 				failReport.FailChecks++
// 			}
// 		}

// 	}

// 	return failReport, errorReport, checkOK
// }

func evaluateOSSExportControl(contributors []common.ContributorMD, controlOpts *gaugeControl, pkgName string) (common.ExportControlSummary, common.FailReport) {
	failReport := common.FailReport{}

	var normalizedBase uint64
	for _, c := range contributors {
		normalizedBase += uint64(c.LOCAdditions)
		normalizedBase += uint64(c.LOCDeletetions)
	}
	if normalizedBase == 0 {
		fmt.Printf("\n** contributor stats not available **\n")
	}
	summary := common.ExportControlSummary{}
	summary.ContributionThreshold = 50
	summary.PackageName = pkgName
	for _, rloc := range controlOpts.ExportControl.RestrictedCountries {
		share := 0
		total := 0
		for _, cloc := range contributors {
			if strings.Contains(strings.ToLower(cloc.ResolvedCountry), strings.ToLower(rloc)) {
				cp := 0
				if normalizedBase != 0 {
					cp = int(((uint64(cloc.LOCAdditions) + uint64(cloc.LOCDeletetions)) * 100 / normalizedBase))
				}
				failReport.Reasons = append(failReport.Reasons, fmt.Sprintf("export_control_location check failed for contributor: `%s`,  location: `%s`, resolved country: `%s`, percent contribution: `%v`", cloc.LoginID, cloc.Location, cloc.ResolvedCountry, cp))
				failReport.NumberOfFailures++
				share += cp
				total++
			}
		}
		summary.Countries = append(summary.Countries, struct {
			CountryName          string "json:\"country_name\""
			Contributors         int    "json:\"total_num_of_contributors\""
			PercentContributions int    "json:\"percent_contributions\""
		}{
			CountryName:          rloc,
			Contributors:         total,
			PercentContributions: share,
		})
	}

	return summary, failReport
}

func evaluatePackageRelease(releaseInsights *common.ReleaseInsights, controlOpts *gaugeControl) common.RecommendedVersion {
	report := common.RecommendedVersion{}

	report.Version = releaseInsights.LatestVersion
	report.ReleaseTimestamp = releaseInsights.LatestReleaseTimestamp
	report.ChangeLabels = releaseInsights.Labels
	report.ZombieChanges = releaseInsights.ChangeInsights.ZombieChanges

	authors := map[string]struct{}{}
	reviewers := map[string]struct{}{}
	var nopeerReview, uniqAuthors, uniqReviewers int
	labels := map[string]struct{}{}
	for _, pr := range releaseInsights.ChangeInsights.Changes {
		if _, found := authors[pr.Authors]; !found {
			uniqAuthors++
			authors[pr.Authors] = struct{}{}
		}

		for _, r := range pr.Approvers {
			if _, found := reviewers[r.LoginName]; !found {
				uniqReviewers++
				reviewers[r.LoginName] = struct{}{}
			}
			if strings.Compare(pr.Authors, r.LoginName) == 0 {
				nopeerReview++
			}
		}
		for _, l := range pr.Labels {
			if _, found := labels[l]; !found {
				labels[l] = struct{}{}
			}
		}
	}
	report.NonPeerReviewedPRs = nopeerReview
	report.NumAuthors = uniqAuthors
	report.NumReviewers = uniqReviewers
	for k := range labels {
		report.ChangeLabels = append(report.ChangeLabels, k)
	}
	return report
}
