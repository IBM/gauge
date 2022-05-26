package ghapis

import (
	"context"
	"net/http"

	"github.com/IBM/gauge/pkg/common"
	"github.com/pkg/errors"
)

// //GetRepositoryMeta :
// func GetRepositoryMeta(ctx context.Context, pkg *common.PackageProps) error {
// 	oauth := os.Getenv(common.GITHUB_API_KEY)
// 	ts := oauth2.StaticTokenSource(
// 		&oauth2.Token{AccessToken: oauth},
// 	)
// 	tc := oauth2.NewClient(ctx, ts)
// 	owner := ""
// 	if pkg.SourceRepo != "" {
// 		owner = parseRepositoryOwner(pkg.SourceRepo)
// 	}
// 	repo := parseRepositoryName(pkg.SourceRepo)
// 	client := github.NewClient(tc)
// 	r, _, err := client.Repositories.Get(ctx, owner, repo)
// 	if err != nil {
// 		errors.Wrapf(err, "error fetching repository metadata")
// 	}
// 	// fmt.Printf("ghclient status code: %d\n", res.StatusCode)
// 	if r != nil {
// 		if r.GetLicense() != nil {
// 			pkg.License = *r.GetLicense().Name
// 		}
// 	}
// 	// for _, r := range rlist {
// 	// 	fmt.Printf("release %s %v %d %s\n", r.GetAssetsURL(), r.GetCreatedAt(), r.GetID(), r.GetTagName())
// 	// }
// 	return nil
// }

//GetRepositoryMD :
func (ghcli *GHClient) GetRepositoryMD(ctx context.Context, pkgName, ecosystem, repo string, result *common.PackageRepoMD) error {
	result.RepoURL = repo
	result.Ecosystem = ecosystem
	result.PackageName = pkgName
	owner := ""
	if repo != "" {
		owner = parseRepositoryOwner(repo)
	}
	repoName := parseRepositoryName(repo)
	r, resp, err := ghcli.ClientV3.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		// fmt.Println(err)
		if resp.StatusCode == http.StatusForbidden {
			return &common.GithubRateLimitErr{}
		}
		return errors.Wrapf(err, "error fetching repository metadata")
	}
	// fmt.Printf("ghclient status code: %d\n", res.StatusCode)
	if r != nil {
		if r.GetLicense() != nil {
			result.License = *r.GetLicense().Name
		}
	}
	result.Size = r.GetSize()
	result.LastUpdated = r.GetUpdatedAt().Time
	// for _, r := range rlist {
	// 	fmt.Printf("release %s %v %d %s\n", r.GetAssetsURL(), r.GetCreatedAt(), r.GetID(), r.GetTagName())
	// }
	return nil
}
