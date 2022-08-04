package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/IBM/gauge/pkg/common"
	"github.com/IBM/gauge/pkg/core"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/pkg/errors"
)

const (
	defaultGaugeConfigFile = ".gauge.yaml"
)

//Package :
func Package() *ffcli.Command {
	var (
		flagset        = flag.NewFlagSet("package", flag.ExitOnError)
		pkgname        = flagset.String("p", "", "package name")
		ecosystem      = flagset.String("e", "", "package ecosystem (node, python, maven,..)")
		releasetag     = flagset.String("t", "", "package release identifier (default: latest)")
		basereleasetag = flagset.String("b", "", "base release to compare against (optional: default to previous)")
		repoURL        = flagset.String("r", "", "repository URL")
		configfp       = flagset.String("c", "", "configuration file")
		outputfp       = flagset.String("f", "", "result filepath")
		deepscan       = flagset.Bool("d", false, "enable deep scan (could be blocked by github API rate-limit)")
	)
	return &ffcli.Command{
		Name:       "package",
		ShortUsage: "gauge package -p <pkgname> -e <ecosystem> -t <release-tag> -r <repo-url> -c <config file> -f <result file>",
		ShortHelp:  `gauge package help`,
		LongHelp: `gauge package help 
EXAMPLES
  # unpack release metadata
  gauge package -p flask -e python -t 2.0.2 -r https://github.com/pallets/flask -c criteria.yaml
`,
		FlagSet: flagset,
		Exec: func(ctx context.Context, args []string) error {

			if *pkgname == "" {
				fmt.Errorf("missing input parameters")
				return errors.New("missing params")
			}
			if *ecosystem == "" && *repoURL == "" {
				fmt.Errorf("missing input parameters")
				return errors.New("missing params")
			}
			if os.Getenv(common.GITHUB_API_KEY) == "" {
				fmt.Errorf("please set `GITHUB_API_KEY` env variable")
				return errors.New("missing params")
			}
			gopts := common.GaugeOpts{}
			gopts.PkgName = *pkgname
			gopts.Ecosystem = *ecosystem
			gopts.ReleaseID = *releasetag
			gopts.RepoURL = *repoURL
			gopts.BaseReleaseID = *basereleasetag
			gopts.PackageOptSelected = true
			gopts.ControlFilepath = *configfp
			gopts.ResultFilepath = *outputfp
			gopts.DeepScanEnabled = *deepscan
			if gopts.ControlFilepath == "" {
				pwd, _ := os.Getwd()
				gopts.ControlFilepath = path.Join(pwd, defaultGaugeConfigFile)
			}
			if err := UnpackRelease(ctx, gopts); err != nil {
				return errors.Wrapf(err, "unpack task for failed")
			}
			return nil
		},
	}
}

//UnpackRelease :
func UnpackRelease(ctx context.Context, dopts common.GaugeOpts) error {
	core.Start(ctx, dopts)
	return nil
}
