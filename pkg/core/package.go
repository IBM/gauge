package core

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/gauge/pkg/common"
	"github.com/IBM/gauge/pkg/ghapis"
	"github.com/IBM/gauge/pkg/releaselib"
	"github.com/IBM/gauge/pkg/utils"
)

func getPkgOSSMeta(ctx context.Context, ghcli *ghapis.GHClient, pkgName, ecosystem, repo string) (common.PackageRepoMD, error) {
	result := common.PackageRepoMD{}
	result.RepoURL = repo

	var cacheState, currState common.PackageRepoMD
	printStatus(ctx, "Going to fetch the cache state...")
	err := releaselib.GetCacheState(pkgName, ecosystem, repo, &cacheState)
	if err != nil {
		printStatus(ctx, "\rGoing to fetch the cache state... \u274c")
	} else {
		printStatus(ctx, "\rGoing to fetch the cache state... \u2714")
	}
	printStatus(ctx, "\n")
	// if repo == "" {
	// 	// fmt.Printf("finding repo url for package: %s\n", pkg.Name)
	// 	repourl := ""
	// 	if ecosystem == common.NODE_ECOSYSTEM {
	// 		repourl, _ = pkgmgr.FetchGitRepositoryFromNPM(pkgName)
	// 	} else if ecosystem == common.PYTHON_ECOSYSTEM {
	// 		repourl, _ = pkgmgr.FetchGitRepositoryFromPYPI(pkgName)
	// 	}
	// 	if repourl != "" {
	// 		repo = repourl
	// 		result.RepoURL = repourl
	// 	}
	// }
	if repo == "" {
		return result, fmt.Errorf("error resolving repository url")
	}

	printStatus(ctx, "Going to fetch the current state from github...")
	errCurrState := ghcli.GetRepositoryMD(ctx, pkgName, ecosystem, repo, &currState)
	if errCurrState != nil {
		printStatus(ctx, "\rGoing to fetch the current state from github... \u274c")
	} else {
		printStatus(ctx, "\rGoing to fetch the current state from github... \u2714")
	}
	printStatus(ctx, "\n")
	// if errCurrState == nil && currState.LastUpdated.Equal(cacheState.LastUpdated) {
	//following change was made to ensure we do not trash for the repos that are getting updated frequently
	if errCurrState == nil && time.Duration(currState.LastUpdated.Sub(cacheState.LastUpdated).Hours()) <= time.Duration(24) {
		printStatus(ctx, "Using cache state for reporting... \u2714")
		result = cacheState
	} else {
		printStatus(ctx, "Getting latest metadata from github... \u2714")
		fmt.Println()
		err := ghcli.GetCodeContributorsMD(ctx, &currState)
		if err != nil {
			printStatus(ctx, "\rGetting latest metadata from github... \u274c")
			//in case of error, return cache state
			return cacheState, err
		}
		result = currState
		printStatus(ctx, "Updating cache state...")
		if err := releaselib.UpdateCacheState(&cacheState, &currState); err != nil {
			printStatus(ctx, "\rUpdating cache state... \u274c")
		}
		printStatus(ctx, "\rUpdating cache state... \u2714")
		printStatus(ctx, "\n")
	}
	printStatus(ctx, "\n")
	if pkgName != "" {
		result.PackageName = pkgName
	}
	if ecosystem != "" {
		result.Ecosystem = ecosystem
	}
	return result, nil
}

func printStatus(ctx context.Context, msg string) {
	if ctx.Value(common.CTX_MODE) == common.CTX_PACKAGE {
		fmt.Printf(msg)
	}
}

//GetPackageReleasesMD :
func GetPackageReleasesMD(ctx context.Context, ghclient *ghapis.GHClient, pkgName, ecosystem, repo, version string, deepscan bool) (common.ReleaseInsights, error) {
	releaseReport := common.ReleaseInsights{}

	releaseReport.PackageName = pkgName
	releaseReport.CurrentVersion = version
	releaseReport.Repo = repo

	latestRelease, err := ghclient.GetLatestRelease(ctx, repo)
	if err != nil {
		fmt.Printf("\nfailed to query github apis: %v\n", err)
		return releaseReport, err
	}
	releaseReport.LatestVersion = latestRelease.Tag
	releaseReport.LatestReleaseTimestamp = latestRelease.CreatedAt
	if utils.IsEqual(latestRelease.Tag, version) {
		releaseReport.IsLatest = true
		releaseReport.LatestReleaseTimestamp = latestRelease.CreatedAt
	} else {
		releaseList, err := ghclient.GetAllReleases(ctx, repo)
		if err != nil {
			fmt.Println("\nfailed to query github apis ")
			return releaseReport, err
		}
		currRelease := common.ReleaseMD{}
		for _, r := range releaseList {
			if utils.IsEqual(r.Tag, version) {
				currRelease = r
				break
			}
		}
		releaseReport.MajorReleaseLag = measureMajorReleaseLag(version, releaseList)
		releaseReport.ReleaseLag = measureReleaseLag(version, releaseList)
		releaseReport.ReleaseTimeLag = fmt.Sprintf("%d days", int((latestRelease.CreatedAt.Sub(currRelease.CreatedAt).Hours() / 24)))
	}

	if deepscan {
		chInsights, err := ghclient.GetChangeInsights(ctx, releaseReport.LatestVersion, releaseReport.CurrentVersion, repo)
		if err != nil {
			fmt.Printf("change insights not available\n")
		}
		releaseReport.ChangeInsights = chInsights
	}
	// json, _ := json.MarshalIndent(releaseReport, "", "    ")
	// fmt.Println(string(json))
	return releaseReport, nil
}

func measureReleaseLag(version string, releases []common.ReleaseMD) int {
	lag := 0
	for _, r := range releases {
		if utils.IsGreater(r.Tag, version) {
			lag++
		}
	}
	return lag
}

func measureMajorReleaseLag(version string, releases []common.ReleaseMD) int {
	lag := 0
	for _, r := range releases {
		if utils.IsGreaterMajor(r.Tag, version) {
			lag++
		}
	}
	return lag
}

func getReleaseMD(version string, releases []common.ReleaseMD) common.ReleaseMD {
	for _, r := range releases {
		if utils.IsEqual(r.Tag, version) {
			return r
		}
	}
	return common.ReleaseMD{}
}
