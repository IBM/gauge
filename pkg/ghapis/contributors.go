package ghapis

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/IBM/gauge/pkg/common"
	"github.com/google/go-github/github"
)

// //GetCodeContributors :
// func GetCodeContributors(ctx context.Context, repoURL string) common.PackageContribResult {
// 	contributors := []common.Contributor{}
// 	result := common.PackageContribResult{}
// 	oauth := os.Getenv(common.GITHUB_API_KEY)
// 	ts := oauth2.StaticTokenSource(
// 		&oauth2.Token{AccessToken: oauth},
// 	)
// 	tc := oauth2.NewClient(ctx, ts)
// 	owner := ""
// 	if repoURL != "" {
// 		owner = parseRepositoryOwner(repoURL)
// 	}
// 	repo := parseRepositoryName(repoURL)
// 	client := github.NewClient(tc)
// 	pageCount := 0
// 	var errorCount, anonCount, totalCount int
// 	for {
// 		pageCount++
// 		users, resp, err := client.Repositories.ListContributors(ctx, owner, repo,
// 			&github.ListContributorsOptions{Anon: "1",
// 				ListOptions: github.ListOptions{Page: pageCount, PerPage: 100}})
// 		if err != nil || resp.StatusCode != 200 || len(users) == 0 {
// 			// fmt.Println(err)
// 			break
// 		}
// 		for _, c := range users {
// 			totalCount++
// 			fmt.Printf("%d\n", c.GetContributions())
// 			if c.GetLogin() != "" {
// 				u, _, err := client.Users.Get(ctx, c.GetLogin())
// 				if err != nil {
// 					errorCount++
// 				}
// 				cadd := common.Contributor{}
// 				if u.GetLocation() != "" {
// 					cadd.Location = *u.Location
// 				}
// 				if u.GetCompany() != "" {
// 					cadd.Affiliation = *u.Company
// 				}
// 				cadd.ID = u.GetLogin()
// 				if u.GetName() != "" {
// 					cadd.Name = u.GetName()
// 				}
// 				if c.GetURL() != "" {
// 					cadd.URL = c.GetURL()
// 				}
// 				contributors = append(contributors, cadd)
// 			} else {
// 				anonCount++
// 			}
// 		}
// 	}
// 	result.TotalContributors = totalCount
// 	result.AnonymizedContributors = anonCount
// 	result.ErroredContributors = errorCount
// 	result.Contributors = contributors
// 	result.SourceRepo = repoURL
// 	return result
// }

//GetCodeContributorsMD :
func (ghcli *GHClient) GetCodeContributorsMD(ctx context.Context, pkgResult *common.PackageRepoMD) error {
	// contributors := []common.Contributor{}
	contributors := []common.ContributorMD{}
	repoURL := pkgResult.RepoURL
	owner := ""
	if repoURL != "" {
		owner = parseRepositoryOwner(repoURL)
	}
	repo := parseRepositoryName(repoURL)
	retries := 3
	var cs []*github.ContributorStats
	var err error
	var resp *github.Response
	incomplete := false
	for retries > 0 {
		//To account for github caching when results are not ready
		// https://docs.github.com/en/rest/reference/metrics#statistics
		cs, resp, err = ghcli.ClientV3.Repositories.ListContributorsStats(ctx, owner, repo)
		if resp.StatusCode == http.StatusAccepted {
			retries--
			time.Sleep(500 * time.Millisecond)
			continue
		}
		retries = 0
		if err != nil || resp.StatusCode != 200 || len(cs) == 0 {
			if resp.StatusCode == http.StatusForbidden {
				fmt.Println("github-api limit reached for ListContributorsStats() call", err)
				return &common.GithubRateLimitErr{}
			}
		}
	}
	cMap := map[string]common.ContributorMD{}
	for _, c := range cs {
		currC := common.ContributorMD{}
		currC.LoginID = c.Author.GetLogin()
		currC.Commits = c.GetTotal()
		for _, w := range c.Weeks {
			currC.LOCAdditions += *w.Additions
			currC.LOCDeletetions += *w.Deletions
		}
		cMap[currC.LoginID] = currC
	}

	pageCount := 0
	var errorCount, anonCount, totalCount int

	for {
		pageCount++
		users, resp, err := ghcli.ClientV3.Repositories.ListContributors(ctx, owner, repo,
			&github.ListContributorsOptions{Anon: "1",
				ListOptions: github.ListOptions{Page: pageCount, PerPage: 100}})
		if err != nil || resp.StatusCode != 200 || len(users) == 0 {
			if resp.StatusCode == http.StatusForbidden {
				fmt.Println("github-api limit reached for ListContributors() call", err)
				incomplete = true
			}
			break
		}
		for _, c := range users {
			totalCount++
			if c.GetLogin() != "" {
				u, _, err := ghcli.ClientV3.Users.Get(ctx, c.GetLogin())
				if err != nil {
					errorCount++
					continue
				}
				cadd := common.ContributorMD{}
				if f, ok := cMap[c.GetLogin()]; ok {
					cadd = f
				}
				cadd.LoginID = c.GetLogin()
				if u.GetLocation() != "" {
					cadd.Location = *u.Location
				}
				if u.GetCompany() != "" {
					cadd.Affiliation = u.GetCompany()
				}
				if u.GetName() != "" {
					cadd.Name = u.GetName()
				}
				contributors = append(contributors, cadd)
			} else {
				anonCount++
			}
		}
	}
	pkgResult.TotalContributors = totalCount
	pkgResult.AnonymizedContributors = anonCount
	pkgResult.ErroredContributors = errorCount
	pkgResult.Contributors = contributors
	if incomplete {
		return &common.GithubRateLimitErr{}
	}
	return nil
}
