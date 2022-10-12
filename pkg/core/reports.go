package core

import (
	"fmt"
	"strings"

	"github.com/IBM/gauge/pkg/common"
)

func printGaugeReportForPackage(report common.GaugeReport, opts common.GaugeOpts) {

	fmt.Printf(strings.Repeat("*", 80))
	fmt.Printf("\nGauge Report for package `%s`\n", report.ReleaseReport.Insights.PackageName)
	fmt.Printf(strings.Repeat("*", 80))

	if opts.PackageReleaseEnabled {
		fmt.Printf("\nRelease Measures: ")
		fmt.Printf("\n\tCurrent version: %s", report.ReleaseReport.Insights.CurrentVersion)
		fmt.Printf("\n\tLatest version: %s", report.ReleaseReport.Insights.LatestVersion)
		fmt.Printf("\n\tRelease lag (versions): %d", report.ReleaseReport.Insights.ReleaseLag)
		fmt.Printf("\n\tRelease lag (days): %s\n", report.ReleaseReport.Insights.ReleaseTimeLag)

		fmt.Printf(strings.Repeat("-", 80))
		fmt.Printf("\n\t\tRecommended update ")
		fmt.Printf("\n\t\t Version - %s", report.ReleaseReport.Recommendation.Version)
		fmt.Printf("\n\t\t Release Time - %s", report.ReleaseReport.Recommendation.ReleaseTimestamp)
		fmt.Printf("\n\t\t Num of unique contributors - %d", report.ReleaseReport.Recommendation.NumAuthors)
		fmt.Printf("\n\t\t Num of unique reviewers - %d", report.ReleaseReport.Recommendation.NumReviewers)
		fmt.Printf("\n\t\t Non peer reviewed changes - %d", report.ReleaseReport.Recommendation.NonPeerReviewedPRs)
		fmt.Printf("\n\t\t Num of zombie commits - %d", report.ReleaseReport.Recommendation.ZombieChanges)
		// lStrs, _ := json.Marshal(recommendation.ChangeLabels)
		// fmt.Printf("\n%v\n%v\n", releaseInsights.Labels, recommendation.ChangeLabels)
		fmt.Printf("\n\t\t Change annotations - %s", "['"+strings.Join(report.ReleaseReport.Recommendation.ChangeLabels, `','`)+`']`)
		fmt.Println()
		fmt.Printf(strings.Repeat("-", 80))
		fmt.Println()
	}
	// if opts.ExportControlEnabled {
	// 	fmt.Printf("Export Control Measures: ")
	// 	fmt.Printf("\n\tTotal Num of Contributors: %d", report.ExportControlReport.OSSMeta.TotalContributors)
	// 	fmt.Printf("\n\tAnonymized Contributors: %d", report.ExportControlReport.OSSMeta.AnonymizedContributors)
	// 	fmt.Printf("\n\tContributors failed to query: %d", report.ExportControlReport.OSSMeta.ErroredContributors)
	// 	fmt.Println()
	// 	for _, f := range report.ExportControlReport.ExportControl.Countries {
	// 		fmt.Printf(strings.Repeat("-", 80))
	// 		fmt.Printf("\n\tRestricted Country: %s", f.CountryName)
	// 		fmt.Printf("\n\t\t Number of contributors: %d", f.Contributors)
	// 		fmt.Printf("\n\t\t Percentage of contributors: %d", f.PercentContributions)
	// 		fmt.Println()
	// 		// fmt.Printf(strings.Repeat("-", 80))
	// 	}
	// 	fmt.Println()
	// 	fmt.Printf(strings.Repeat("*", 80))
	// 	fmt.Println()
	// }
}
