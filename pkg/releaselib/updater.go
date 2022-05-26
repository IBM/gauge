package releaselib

import (
	"fmt"
	"os"
	"path"

	"github.com/IBM/gauge/pkg/common"
)

//UpdateCacheState :
func UpdateCacheState(cacheState, currState *common.PackageRepoMD) error {
	_, cacheReady := os.LookupEnv(common.RELEASE_LIB_SERVER)
	if !cacheReady {
		return fmt.Errorf("cache server not setup")
	}

	// check if repository needs to be updated
	if cacheState.LastUpdated != currState.LastUpdated {
		repoMD := common.RepositoryMD{}
		repoMD.License = currState.License
		repoMD.URL = currState.RepoURL
		repoMD.Size = currState.Size
		repoMD.Name = path.Base(currState.RepoURL)
		repoMD.UpdatedAt = currState.LastUpdated
		if err := StoreRepositoryMeta(&repoMD); err != nil {
			return fmt.Errorf("failed to update cache for repository")
		}
		if err := StoreContributors(currState); err != nil {
			return fmt.Errorf("failed to update cache for contributors")
		}
	}

	// check if need to update package info
	if cacheState.PackageName == "" && currState.PackageName != "" {
		pkgMD := common.PackageMD{}
		pkgMD.Ecosystem = currState.Ecosystem
		pkgMD.PackageName = currState.PackageName
		pkgMD.Version = currState.Version
		pkgMD.RepoURL = currState.RepoURL
		if err := StorePackageMeta(&pkgMD); err != nil {
			return fmt.Errorf("failed to update cache for package")
		}
	}

	return nil
}

func diffContributors(cacheState *[]common.ContributorMD, currState *[]common.ContributorMD) []common.ContributorMD {
	diffContribs := []common.ContributorMD{}

	return diffContribs
}
