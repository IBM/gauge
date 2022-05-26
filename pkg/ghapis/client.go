package ghapis

import (
	"context"
	"os"

	"github.com/IBM/gauge/pkg/common"
	"github.com/google/go-github/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

//GHClient :
type GHClient struct {
	ClientV3 *github.Client
	ClientV4 *githubv4.Client
}

//Setup : setup github client for v3 and v4
func (cli *GHClient) Setup(ctx context.Context) error {
	//setup v2 client for go-client apis
	oauth := os.Getenv(common.GITHUB_API_KEY)
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: oauth},
	)
	tc := oauth2.NewClient(ctx, ts)
	cli.ClientV3 = github.NewClient(tc)

	cli.ClientV4 = githubv4.NewClient(tc)

	return nil
}
