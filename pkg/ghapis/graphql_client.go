package ghapis

import (
	"context"
	"fmt"

	"github.com/IBM/gauge/pkg/common"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
)

func (ghcli *GHClient) getAssociatedPRNumber(ctx context.Context, owner, repo, sha string) (int, error) {
	prNum := 0
	var q struct {
		Repository struct {
			Object struct {
				Commit struct {
					AssociatedPullRequests struct {
						Edges []struct {
							Node struct {
								Number int
							}
						}
					} `graphql:"associatedPullRequests(first: 1)"`
				} `graphql:"... on Commit"`
			} `graphql:"object(expression: $sha)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}
	v := map[string]interface{}{
		"owner": githubv4.String(owner),
		"repo":  githubv4.String(repo),
		"sha":   githubv4.String(sha),
	}
	if err := ghcli.ClientV4.Query(ctx, &q, v); err != nil {
		return prNum, errors.Wrapf(err, "error while query")
	}
	if len(q.Repository.Object.Commit.AssociatedPullRequests.Edges) != 0 {
		prNum = q.Repository.Object.Commit.AssociatedPullRequests.Edges[0].Node.Number
	}
	return prNum, nil
}

func (ghcli *GHClient) getPRApprovers(ctx context.Context, owner, repo string, prnum int) ([]common.Approver, error) {
	var q struct {
		Repository struct {
			PullRequest struct {
				Reviews struct {
					Nodes []struct {
						Author struct {
							Login        string
							ResourcePath string
							URL          string
						}
					}
				} `graphql:"reviews(first: 100, states: APPROVED)"`
			} `graphql:"pullRequest(number: $num)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}
	v := map[string]interface{}{
		"owner": githubv4.String(owner),
		"repo":  githubv4.String(repo),
		"num":   githubv4.Int(prnum),
	}
	if err := ghcli.ClientV4.Query(ctx, &q, v); err != nil {
		fmt.Println(err)
		return nil, errors.Wrapf(err, "error while query")
	}
	nodes := q.Repository.PullRequest.Reviews.Nodes
	approvers := []common.Approver{}
	if len(nodes) != 0 {
		for _, i := range nodes {
			approver := common.Approver{
				LoginName:    i.Author.Login,
				ResourcePath: i.Author.ResourcePath,
				URL:          i.Author.URL,
			}
			approvers = append(approvers, approver)
		}
	}
	return approvers, nil
}

func (ghcli *GHClient) getPRLinkedIssues(ctx context.Context, owner, repo string, prnum int) ([]int, error) {
	var q struct {
		Repository struct {
			PullRequest struct {
				ClosingIssueReferences struct {
					Nodes []struct {
						Number int
					}
				} `graphql:"closingIssuesReferences(first: 10)"`
			} `graphql:"pullRequest(number: $num)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}
	v := map[string]interface{}{
		"owner": githubv4.String(owner),
		"repo":  githubv4.String(repo),
		"num":   githubv4.Int(prnum),
	}
	// v := map[string]interface{}{}
	if err := ghcli.ClientV4.Query(ctx, &q, v); err != nil {
		fmt.Println(err)
		return nil, errors.Wrapf(err, "error while query")
	}
	nodes := q.Repository.PullRequest.ClosingIssueReferences.Nodes
	issues := []int{}
	if len(nodes) != 0 {
		for _, i := range nodes {
			issues = append(issues, i.Number)
		}
	}
	return issues, nil
}
