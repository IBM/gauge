package core

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/IBM/gauge/pkg/common"
	"github.com/IBM/gauge/pkg/ghapis"
	"github.com/IBM/gauge/pkg/pkgmgr"
	"github.com/IBM/gauge/pkg/sbom"
	"github.com/cheggaaa/pb/v3"
)

//Start : entry point for core logic
func Start(ctx context.Context, opts common.GaugeOpts) error {

	gaugeCtr := gaugeControl{}
	parseGaugeControls(opts.ControlFilepath, &gaugeCtr)

	ghclient := ghapis.GHClient{}
	ghclient.Setup(ctx)

	if opts.PackageOptSelected {
		runGaugeForPackage(ctx, opts, &gaugeCtr, &ghclient)
	} else if opts.SBOMOptSelected {
		runGaugeForSBOM(ctx, opts, &gaugeCtr, &ghclient)
	}
	return nil
}

func runGaugeForPackage(ctx context.Context, opts common.GaugeOpts, gaugeCtr *gaugeControl, ghclient *ghapis.GHClient) {
	report := common.GaugeReport{}
	opts.RepoURL = resolvePackageSource(opts.PkgName, opts.Ecosystem, opts.RepoURL)
	if gaugeCtr.ReleaseConfig.Enable {
		gaugePackageRelease(ctx, opts, gaugeCtr, ghclient, &report)
	}
	if gaugeCtr.ExportControl.Enable {
		packageGauger(ctx, opts, gaugeCtr, ghclient, &report)
	}
	opts.ExportControlEnabled = gaugeCtr.ExportControl.Enable
	opts.PackageReleaseEnabled = gaugeCtr.ReleaseConfig.Enable

	printGaugeReportForPackage(report, opts)

	if opts.ResultFilepath != "" {
		storeLogReport(report, opts.ResultFilepath)
	}
}

func runGaugeForSBOM(ctx context.Context, opts common.GaugeOpts, gaugeCtr *gaugeControl, ghclient *ghapis.GHClient) error {
	report := []common.GaugeReport{}
	pkgList := []common.PackageProps{}
	if strings.Compare(strings.ToLower(common.CYCLONEDX_VAR), opts.SBOMFormat) == 0 {
		pkgList, _ = sbom.ParseOSSPackagesCyclonedx(opts.SBOMFilepath)
	} else if strings.Compare(strings.ToLower(common.CSV_VAR), opts.SBOMFormat) == 0 {
		pkgList, _ = sbom.ParseOSSPackagesCSV(opts.SBOMFilepath)
	} else {
		return fmt.Errorf("sbom format `%s` currently not supported", opts.SBOMFormat)
	}
	if len(pkgList) == 0 {
		return errors.New("error reading sbom file")
	}
	nPkgs := len(pkgList)
	progress := pb.StartNew(nPkgs)
	totalPkgs := 0
	resolveFailedPkgs := 0
	// cntrFailedPkgs := 0
	ghRateLimitHit := false

	for idx := 0; idx < nPkgs && !ghRateLimitHit; idx++ {
		pkg := pkgList[idx]
		totalPkgs++
		if totalPkgs%5 == 0 {
			progress.SetCurrent(int64(totalPkgs))
		}
		currGaugeReport := common.GaugeReport{}
		opts.RepoURL = resolvePackageSource(pkg.Name, pkg.Ecosystem, "")
		opts.PkgName = pkg.Name
		opts.Ecosystem = pkg.Ecosystem
		opts.ReleaseID = pkg.Key

		var err error
		if gaugeCtr.ReleaseConfig.Enable {
			err = gaugePackageRelease(ctx, opts, gaugeCtr, ghclient, &currGaugeReport)
		}
		if gaugeCtr.ExportControl.Enable {
			err = packageGauger(ctx, opts, gaugeCtr, ghclient, &currGaugeReport)
		}
		if err != nil {
			resolveFailedPkgs++
			if errors.Is(err, &common.GithubRateLimitErr{}) {
				fmt.Printf("Github API rate limit reached, following result maybe incomplete. ")
				ghRateLimitHit = true
			}
		}
		report = append(report, currGaugeReport)
		if idx == 5 {
			break
		}
	}

	progress.Finish()
	// printSBOMReport(opts.SBOMFilepath, guageCtr.ExportControl.RestrictedCountries, &summary)
	// fmt.Printf("complete log file is available at: %s\n", logFile)

	if opts.ResultFilepath != "" {
		storeLogReport(report, opts.ResultFilepath)
	}

	return nil
}

func sbomGauger(ctx context.Context, opts common.GaugeOpts, guageCtr *gaugeControl, ghcli *ghapis.GHClient) (common.ExportControlSummarySBOM, error) {
	summary := common.ExportControlSummarySBOM{}
	logreport := []common.LogReport{}

	pkgList := []common.PackageProps{}
	if strings.Compare(strings.ToLower(common.CYCLONEDX_VAR), opts.SBOMFormat) == 0 {
		pkgList, _ = sbom.ParseOSSPackagesCyclonedx(opts.SBOMFilepath)
	} else if strings.Compare(strings.ToLower(common.CSV_VAR), opts.SBOMFormat) == 0 {
		pkgList, _ = sbom.ParseOSSPackagesCSV(opts.SBOMFilepath)
	} else {
		return summary, fmt.Errorf("sbom format `%s` currently not supported", opts.SBOMFormat)
	}
	if len(pkgList) == 0 {
		return summary, errors.New("error reading sbom file")
	}

	nPkgs := len(pkgList)
	progress := pb.StartNew(nPkgs)
	totalPkgs := 0
	resolveFailedPkgs := 0
	// cntrFailedPkgs := 0
	pkgsReport := []common.ExportControlSummary{}
	ghRateLimitHit := false

	for idx := 0; idx < nPkgs && !ghRateLimitHit; idx++ {
		pkg := pkgList[idx]
		totalPkgs++
		pkglogger := common.LogReport{}
		pkgSummary := common.ExportControlSummary{}
		pkglogger.PackageName = pkg.Name
		if totalPkgs%5 == 0 {
			progress.SetCurrent(int64(totalPkgs))
		}
		ossMeta, err := getPkgOSSMeta(ctx, ghcli, pkg.Name, pkg.Ecosystem, pkg.SourceRepo)
		if err != nil {
			pkglogger.PackageErr = common.ErrorReport{
				NumberOfError: 1,
				Reasons:       []string{fmt.Sprintf("package metadata not resolved")},
			}
			resolveFailedPkgs++
			logreport = append(logreport, pkglogger)
			if errors.Is(err, &common.GithubRateLimitErr{}) {
				fmt.Printf("Github API rate limit reached, following result maybe incomplete. ")
				ghRateLimitHit = true
			}
		} else {
			pkglogger.LocationErr = resolveLocation(ossMeta.Contributors)
			pkgSummary, pkglogger.ControlErr = evaluateOSSExportControl(ossMeta.Contributors, guageCtr, pkg.Name)
			pkglogger.RepoURL = ossMeta.RepoURL
			pkglogger.LastUpdated = ossMeta.LastUpdated
			pkgsReport = append(pkgsReport, pkgSummary)
			logreport = append(logreport, pkglogger)
		}
	}
	logFile, _ := storeLogReport(logreport, "")
	summary.AnalyzedPackages = totalPkgs - resolveFailedPkgs
	summary.TotalPackages = nPkgs
	summarizeSBOMResults(pkgsReport, &summary)
	progress.Finish()
	printSBOMReport(opts.SBOMFilepath, guageCtr.ExportControl.RestrictedCountries, &summary)
	fmt.Printf("complete log file is available at: %s\n", logFile)
	return summary, nil
}

func summarizeSBOMResults(pkgRes []common.ExportControlSummary, sbomRes *common.ExportControlSummarySBOM) {
	// contribByCountry := map[string]int{}
	check := common.ControlCheck{}
	failedCnt := 0
	for _, pkg := range pkgRes {
		if len(pkg.Countries) != 0 {
			failedCnt++
			check.ListOfPackages = append(check.ListOfPackages, pkg.PackageName)
			// for _, c := pkg.Countries{
			// 	if v, ok := contribByCountry[c]; ok {
			// 	}
			// }
		}
	}
	sbomRes.CheckFails = []common.ControlCheck{check}
}

func packageGauger(ctx context.Context, opts common.GaugeOpts, guageCtr *gaugeControl, ghcli *ghapis.GHClient, report *common.GaugeReport) error {
	logreport := common.LogReport{}
	summary := common.ExportControlSummary{}

	ossMeta, err := getPkgOSSMeta(ctx, ghcli, opts.PkgName, opts.Ecosystem, opts.RepoURL)
	if err != nil {
		return err
	}

	logreport.LocationErr = resolveLocation(ossMeta.Contributors)
	summary, logreport.ControlErr = evaluateOSSExportControl(ossMeta.Contributors, guageCtr, ossMeta.PackageName)

	logreport.LastUpdated = ossMeta.LastUpdated
	logreport.PackageName = ossMeta.PackageName
	logreport.RepoURL = ossMeta.RepoURL

	// logFile, err := storeLogReport(logreport, "")
	// // printPackageReport(summary, ossMeta)
	// fmt.Printf("complete log file is available at: %s\n", logFile)

	report.ExportControlReport.OSSMeta = ossMeta
	report.ExportControlReport.ExportControl = summary
	report.ExportControlReport.ExportControl.AuditLogs = logreport

	return nil
}

func printSBOMReport(sbomfp string, exportCntrPolicy []string, report *common.ExportControlSummarySBOM) {
	fmt.Printf(strings.Repeat("*", 80))
	fmt.Printf("\nGauge Report for SBOM %s:\n", sbomfp)
	fmt.Printf(strings.Repeat("*", 80))
	fmt.Printf("\nTotal Num of packages in SBOM: %d", report.TotalPackages)
	fmt.Printf("\nPackages analyzed successfully: %d", report.AnalyzedPackages)
	fmt.Printf("\nExport control policy for restricted countries: %v", "'"+strings.Join(exportCntrPolicy, `','`)+`'`)
	fmt.Printf("\nNumber of packages failed: %d", len(report.CheckFails[0].ListOfPackages))
	fmt.Printf("\nList of packages failed: %v", "'"+strings.Join(report.CheckFails[0].ListOfPackages, `','`)+`'`)
	fmt.Println()
	fmt.Printf(strings.Repeat("*", 80))
	fmt.Println()
}

func printPackageReport(report common.ExportControlSummary, ossMeta common.PackageRepoMD) {
	fmt.Printf(strings.Repeat("*", 80))
	fmt.Printf("\nGauge Report for package %s:\n", ossMeta.PackageName)
	fmt.Printf(strings.Repeat("*", 80))
	fmt.Printf("\nTotal Num of Contributors: %d", ossMeta.TotalContributors)
	fmt.Printf("\nAnonymized Contributors: %d", ossMeta.AnonymizedContributors)
	fmt.Printf("\nContributors failed to query: %d", ossMeta.ErroredContributors)
	fmt.Println()
	for _, f := range report.Countries {
		fmt.Printf(strings.Repeat("-", 80))
		fmt.Printf("\nRestricted Country: %s", f.CountryName)
		fmt.Printf("\n\t Number of contributors: %d", f.Contributors)
		fmt.Printf("\n\t Percentage of contributors: %d", f.PercentContributions)
		fmt.Println()
		// fmt.Printf(strings.Repeat("-", 80))
	}
	fmt.Println()
	fmt.Printf(strings.Repeat("*", 80))
	fmt.Println()
}

func storeLogReport(logreport interface{}, filepath string) (string, error) {
	if filepath == "" {
		fp, err := ioutil.TempFile(os.TempDir(), "gauge-")
		if err != nil {
			return "", err
		}
		filepath = fp.Name()
	}
	logf, _ := os.Create(filepath)
	writer := bufio.NewWriter(logf)
	logBuf, _ := json.MarshalIndent(logreport, "", "    ")

	writer.Write(logBuf)

	return logf.Name(), nil
}

func gaugePackageRelease(ctx context.Context, opts common.GaugeOpts, gaugeCtr *gaugeControl, ghclient *ghapis.GHClient, report *common.GaugeReport) error {
	releaseInsights, err := GetPackageReleasesMD(ctx, ghclient, opts.PkgName, opts.Ecosystem, opts.RepoURL, opts.ReleaseID, opts.DeepScanEnabled)
	if err != nil {
		fmt.Printf("failed to query package release information\n")
		return fmt.Errorf("failed to query package release information")
	}

	recommendation := evaluatePackageRelease(&releaseInsights, gaugeCtr)
	// printPackageReleaseReport(releaseInsights, recommendation)
	report.ReleaseReport.Insights = releaseInsights
	report.ReleaseReport.Recommendation = recommendation

	return nil
}

func printPackageReleaseReport(releaseInsights common.ReleaseInsights, recommendation common.RecommendedVersion) {
	fmt.Printf(strings.Repeat("*", 80))
	fmt.Printf("\nGauge Report for package %s:\n", releaseInsights.PackageName)
	fmt.Printf(strings.Repeat("*", 80))

	fmt.Printf("\nCurrent version: %s", releaseInsights.CurrentVersion)
	fmt.Printf("\nLatest version: %s", releaseInsights.LatestVersion)
	fmt.Printf("\nRelease lag (versions): %d", releaseInsights.ReleaseLag)
	fmt.Printf("\nRelease lag (days): %s\n", releaseInsights.ReleaseTimeLag)

	fmt.Printf(strings.Repeat("-", 80))
	fmt.Printf("\nRecommended update ")
	fmt.Printf("\n\t Version - %s", recommendation.Version)
	fmt.Printf("\n\t Release Time - %s", recommendation.ReleaseTimestamp)
	fmt.Printf("\n\t Num of unique contributors - %d", recommendation.NumAuthors)
	fmt.Printf("\n\t Num of unique reviewers - %d", recommendation.NumReviewers)
	fmt.Printf("\n\t Non peer reviewed changes - %d", recommendation.NonPeerReviewedPRs)
	fmt.Printf("\n\t Num of zombie commits - %d", recommendation.ZombieChanges)
	// lStrs, _ := json.Marshal(recommendation.ChangeLabels)
	// fmt.Printf("\n%v\n%v\n", releaseInsights.Labels, recommendation.ChangeLabels)
	fmt.Printf("\n\t Change annotations - %s", "['"+strings.Join(recommendation.ChangeLabels, `','`)+`']`)
	fmt.Println()
	fmt.Printf(strings.Repeat("-", 80))
	fmt.Println()

}

func resolvePackageSource(pkgName, ecosystem, repourl string) string {
	if repourl == "" {
		if ecosystem == common.NODE_ECOSYSTEM {
			repourl, _ = pkgmgr.FetchGitRepositoryFromNPM(pkgName)
		} else if ecosystem == common.PYTHON_ECOSYSTEM {
			repourl, _ = pkgmgr.FetchGitRepositoryFromPYPI(pkgName)
		}
	}
	return repourl
}
