package ghapis

import (
	"context"
)

func (cli *GHClient) getIssueLabels(ctx context.Context, owner, repo string, issues []int) []string {
	labels := []string{}
	for _, issue := range issues {
		i, _, _ := cli.ClientV3.Issues.Get(ctx, owner, repo, issue)
		for _, l := range i.Labels {
			labels = append(labels, l.GetName())
		}
	}
	return labels
}
