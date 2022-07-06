package ghapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/IBM/gauge/pkg/common"
	"github.com/google/go-github/github"
)

//GetAllReleases :
func (gaugeCli *GHClient) GetAllReleases(ctx context.Context, repoURL string) ([]common.ReleaseMD, error) {
	releases := []common.ReleaseMD{}
	repoOwner := ""
	if repoURL != "" {
		repoOwner = parseRepositoryOwner(repoURL)
	}
	repo := parseRepositoryName(repoURL)
	pageCount := 0
	for {
		pageCount++
		rlist, ghres, err := gaugeCli.ClientV3.Repositories.ListReleases(ctx, repoOwner, repo, &github.ListOptions{Page: pageCount, PerPage: 100})
		if err != nil || ghres.StatusCode != 200 || len(rlist) == 0 {
			if ghres.StatusCode == http.StatusForbidden {
				return releases, &common.GithubRateLimitErr{}
			}
			break
		}
		for _, r := range rlist {
			releases = append(releases, common.ReleaseMD{
				Tag:       r.GetTagName(),
				CreatedAt: r.GetCreatedAt().Time,
				CommitID:  r.GetTargetCommitish(),
			})
		}
	}
	return releases, nil
}

//GetLatestRelease :
func (gaugeCli *GHClient) GetLatestRelease(ctx context.Context, repoURL string) (common.ReleaseMD, error) {
	release := common.ReleaseMD{}
	repoOwner := ""
	if repoURL != "" {
		repoOwner = parseRepositoryOwner(repoURL)
	}
	repo := parseRepositoryName(repoURL)
	result, ghresp, err := gaugeCli.ClientV3.Repositories.GetLatestRelease(ctx, repoOwner, repo)
	if err != nil {
		if ghresp.StatusCode == http.StatusForbidden {
			return release, &common.GithubRateLimitErr{}
		}
		return release, errors.Wrapf(err, "error quering releases")
	}
	if ghresp.StatusCode != 200 {
		return release, errors.Wrapf(err, "un-expected response code %d\n", ghresp.StatusCode)
	}
	release.Tag = result.GetTagName()
	release.CreatedAt = result.GetCreatedAt().Time
	release.CommitID = result.GetTargetCommitish()
	return release, nil
}

func (ghcli *GHClient) getReleaseTimestamp(ctx context.Context, owner, repo, releaseID string) (github.Timestamp, error) {
	ts := github.Timestamp{}
	release, ghresp, err := ghcli.ClientV3.Repositories.GetReleaseByTag(ctx, owner, repo, releaseID)
	if err != nil {
		fmt.Printf("", err)
		return ts, errors.Wrapf(err, "error quering releases")
	}
	if ghresp.StatusCode != 200 {
		fmt.Println(ghresp.StatusCode)
		return ts, errors.Wrapf(err, "un-expected response code %d\n", ghresp.StatusCode)
	}
	return release.GetPublishedAt(), nil
}

//GetChangeInsights :
func (ghcli *GHClient) GetChangeInsights(ctx context.Context, releaseID, baseReleaseID, repoURL string) (common.CommitHistory, error) {
	commitH := common.CommitHistory{}
	repoOwner := ""
	if repoURL != "" {
		repoOwner = parseRepositoryOwner(repoURL)
	}
	repo := parseRepositoryName(repoURL)
	l, resp, err := ghcli.ClientV3.Repositories.ListTags(ctx, repoOwner, repo, &github.ListOptions{})
	if err != nil {
		if resp.StatusCode == http.StatusForbidden {
			return commitH, &common.GithubRateLimitErr{}
		}
		return commitH, errors.Wrapf(err, "error quering releases")
	}
	previousReleaseCommit := ""
	currentReleaseCommit := ""
	findNext := false
	for _, i := range l {
		if findNext && baseReleaseID != "" {
			if i.GetName() == baseReleaseID {
				previousReleaseCommit = *i.GetCommit().SHA
				// releaseMeta.BaseReleaseTag = i.GetName()
				// releaseMeta.BaseReleaseTime, _ = ghcli.getReleaseTimestamp(ctx, repoOwner, repo, baseReleaseID)
				break
			}
		}
		if i.GetName() == releaseID {
			currentReleaseCommit = *i.GetCommit().SHA
			// releaseMeta.ReleaseTime, _ = ghcli.getReleaseTimestamp(ctx, repoOwner, repo, releaseID)
			findNext = true
		}
	}
	commitH.Changes = []common.PullRequest{}
	prCache := map[int]struct{}{}
	// fmt.Println(previousReleaseCommit, "--", currentReleaseCommit)
	gc, resp, err := ghcli.ClientV3.Repositories.CompareCommits(ctx, repoOwner, repo, previousReleaseCommit, currentReleaseCommit)
	if err != nil {
		if resp.StatusCode == http.StatusForbidden {
			return commitH, &common.GithubRateLimitErr{}
		}
		return commitH, errors.Wrapf(err, "error quering releases")
	}
	for _, c := range gc.Commits {
		prnum, _ := ghcli.getAssociatedPRNumber(ctx, repoOwner, repo, c.GetSHA())
		if prnum != 0 && !contains(prCache, prnum) {
			pr := ghcli.getPRMeta(ctx, repoOwner, repo, prnum)
			pr.Authors = c.GetAuthor().GetLogin()
			commitH.Changes = append(commitH.Changes, pr)
			prCache[prnum] = struct{}{}
		} else {
			// zombieChange := common.PullRequest{
			// 	Commits: 1,
			// }
			// releaseMeta.Changes = append(releaseMeta.Changes, zombieChange)
			commitH.ZombieChanges++
		}
	}
	// r, _ := json.MarshalIndent(commitH, "", "    ")
	// fmt.Println("release meta : %v", string(r))
	return commitH, nil
}

func (ghcli *GHClient) getPRMeta(ctx context.Context, owner, repo string, prnum int) common.PullRequest {
	prRes := common.PullRequest{}
	probj, resp, err := ghcli.ClientV3.PullRequests.Get(ctx, owner, repo, prnum)
	if resp.StatusCode != http.StatusOK || err != nil {
		return prRes
	}
	prRes.PRurl = probj.GetURL()
	prRes.Labels = []string{}
	prRes.Timestamp = probj.GetClosedAt().String()
	prRes.Commits = probj.GetCommits()
	prRes.IssueURL = probj.GetIssueURL()

	// fmt.Println("author assoc ", probj.GetAuthorAssociation())
	uniqLabels := map[string]struct{}{}
	if probj != nil {
		for _, l := range probj.Labels {
			if _, found := uniqLabels[l.GetName()]; !found {
				prRes.Labels = append(prRes.Labels, l.GetName())
				uniqLabels[l.GetName()] = struct{}{}
			}
		}
		prRes.Approvers, _ = ghcli.getPRApprovers(ctx, owner, repo, prnum)
	}

	issues, err := ghcli.getPRLinkedIssues(ctx, owner, repo, prnum)
	if err == nil && len(issues) != 0 {
		issueLabels := ghcli.getIssueLabels(ctx, owner, repo, issues)
		if len(issueLabels) != 0 {
			for _, l := range issueLabels {
				if _, found := uniqLabels[l]; !found {
					prRes.Labels = append(prRes.Labels, l)
					uniqLabels[l] = struct{}{}
				}
			}
		}
	}
	return prRes
}
