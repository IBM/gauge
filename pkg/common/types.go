package common

import (
	"time"
)

//GaugeOpts :
type GaugeOpts struct {
	UnpackOpts
	SBOMInputOpts
	SBOMOptSelected       bool
	PackageOptSelected    bool
	ExportControlEnabled  bool
	PackageReleaseEnabled bool
	ResultFilepath        string
}

type SBOMInputOpts struct {
	SBOMFilepath    string
	SBOMFormat      string
	ControlFilepath string
	DeepScanEnabled bool
}

//UnpackOpts :
type UnpackOpts struct {
	PkgName       string
	Ecosystem     string
	ReleaseID     string
	RepoURL       string
	BaseReleaseID string
}

const (
	GITHUB_API_KEY      = "GITHUB_API_KEY"
	LIBRARIESIO_API_KEY = "LIBRARIESIO_API_KEY"
	RELEASE_LIB_SERVER  = "RELEASE_LIB_SERVER"
	WEATHER_API_KEY     = "WEATHER_API_KEY"

	CYCLONEDX_VAR = "cycloneDX"
	SPDX_VAR      = "spdx"
	CSV_VAR       = "csv"

	PYTHON_ECOSYSTEM = "Python"
	NODE_ECOSYSTEM   = "JavaScript"
	GO_ECOSYSTEM     = "Go"

	CTX_MODE    = "mode"
	CTX_PACKAGE = "package"
	CTX_SBOM    = "sbom"
)

type Approver struct {
	LoginName    string
	ResourcePath string
	URL          string
}

type PullRequest struct {
	PRurl     string
	Authors   string
	Timestamp string
	Approvers []Approver
	Commits   int
	Labels    []string
	IssueURL  string
}

type CommitHistory struct {
	Changes       []PullRequest
	Contributors  int
	Approvers     int
	ZombieChanges int
}

type PackageProps struct {
	Name       string  `json:"name"`
	SourceRepo string  `json:"source_repo"`
	Scorecard  float64 `json:"scorecard"`
	Ecosystem  string  `json:"ecosystem"`
	License    string  `json:"license"`
	Key        string  `json:"version"`
}

//Contributor :
type Contributor struct {
	Name        string `json:"name"`
	ID          string `json:"login_id"`
	URL         string `json:"url"`
	Affiliation string
	Location    string `json:"location"`
}

// //ProjectContributorResult :
// type ProjectContributorResult struct {
// 	TotalContributors      int           `json:"total_contributors"`
// 	AnonymizedContributors int           `json:"anonymized_contributors"`
// 	ErroredContributors    int           `json:"errored_contributors"`
// 	Contributors           []Contributor `json:"contributors"`
// }

//PackageHealthStat :
type PackageHealthStat struct {
	Name                   string `json:"name"`
	Version                string `json:"version"`
	Contributors           []Contributor
	TotalContributors      int `json:"total_contributors"`
	AnonymizedContributors int `json:"anonymized_contributors"`
	ErroredContributors    int `json:"errored_contributors"`
	ReleaseProvenance      CommitHistory
}

//EvalReport :
type EvalReport struct {
	NumPkgScanned int                  `json:"num_packages_scanned"`
	NumPass       int                  `json:"num_packages_passed"`
	NumFailed     int                  `json:"num_packages_failed"`
	FailReport    []FailePackageReport `json:"failed_package_report"`
	ErrorReport   []ErrorPackageReport `json:"error_package_report"`

	TotalContributors int `json:"total_contributors"`
	AnonContributors  int `json:"anon_contributors"`
	ErrorContributors int `json:"error_contributors"`
}

//EvalPackageReport :
type EvalPackageReport struct {
	PackageName       string      `json:"package_name"`
	PackageVersion    string      `json:"package_version"`
	TotalContributors int         `json:"total_contributors"`
	AnonContributors  int         `json:"anonymized_contributors"`
	ErrContributors   int         `json:"err_contributors"`
	Failedchecks      FailReport  `json:"failed_checks"`
	ErrorChecks       ErrorReport `json:"error_checks"`
}

//FailReport :
type FailReport struct {
	NumberOfFailures int      `json:"number_of_failures"`
	Reasons          []string `json:"reasons"`
}

//ErrorReport :
type ErrorReport struct {
	NumberOfError int      `json:"number_of_errors"`
	Reasons       []string `json:"reasons"`
}

//ErrorPackageReport :
type ErrorPackageReport struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	TotalChecks int
	ErrorChecks int
	Reasons     []string `json:"reasons"`
}

//FailePackageReport :
type FailePackageReport struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	TotalChecks int
	FailChecks  int
	Reasons     []string `json:"reasons"`
}

//PackageContribResult :
type PackageContribResult struct {
	PackageName            string        `json:"package_name"`
	SourceRepo             string        `json:"source_repo"`
	Scorecard              float64       `json:"scorecard"`
	Ecosystem              string        `json:"ecosystem"`
	License                string        `json:"license"`
	Contributors           []Contributor `json:"contributors"`
	TotalContributors      int           `json:"total_contributors"`
	AnonymizedContributors int           `json:"anonymized_contributors"`
	ErroredContributors    int           `json:"errored_contributors"`
}

//PackageRepoMD :
type PackageRepoMD struct {
	PackageName            string          `json:"package_name"`
	RepoURL                string          `json:"repo_url"`
	Version                string          `json:"version"`
	LastUpdated            time.Time       `json:"updated_at"`
	TotalContributors      int             `json:"total_contributors"`
	Size                   int             `json:"size"`
	License                string          `json:"license"`
	Ecosystem              string          `json:"ecosystem"`
	Contributors           []ContributorMD `json:"-"`
	AnonymizedContributors int             `json:"anonymized_contributors"`
	ErroredContributors    int             `json:"errored_contributors"`
}

//ContributorMD :
type ContributorMD struct {
	Name            string `json:"name"`
	LoginID         string `json:"login_id"`
	Location        string `json:"location"`
	Affiliation     string `json:"affiliation"`
	ResolvedCountry string `json:"resolved_country"`
	LOCAdditions    int    `json:"additions"`
	LOCDeletetions  int    `json:"deletion"`
	Commits         int    `json:"commits"`
}

//RepositoryMD :
type RepositoryMD struct {
	Name      string    `json:"name"`
	License   string    `json:"license"`
	UpdatedAt time.Time `json:"updated_at"`
	URL       string    `json:"repo_url"`
	Size      int
}

//PackageMD :
type PackageMD struct {
	PackageName string    `json:"package_name"`
	RepoURL     string    `json:"repo_url"`
	Version     string    `json:"package_version"`
	LastUpdated time.Time `json:"updated_at"`
	Ecosystem   string    `json:"ecosystem"`
	License     string    `json:"license"`
	Scorecard   float64   `json:"scorecard_score"`
}

//ExportControlSummary :
type ExportControlSummary struct {
	PackageName           string `json:"package_name"`
	ContributionThreshold int
	Countries             []struct {
		CountryName          string `json:"country_name"`
		Contributors         int    `json:"total_num_of_contributors"`
		PercentContributions int    `json:"percent_contributions"`
	}
	AuditLogs LogReport `json:"audit_log_report"`
}

//LogReport :
type LogReport struct {
	PackageName string      `json:"package_name"`
	RepoURL     string      `json:"repo_url"`
	LastUpdated time.Time   `json:"updated_at"`
	LocationErr ErrorReport `json:"location_check_failed"`
	PackageErr  ErrorReport `json:"package_resolve_failed"`
	ControlErr  FailReport  `json:"control_check_failed"`
}

//ExportControlSummarySBOM :
type ExportControlSummarySBOM struct {
	TotalPackages    int            `json:"total_packages"`
	AnalyzedPackages int            `json:"analyzed_packages"`
	CheckFails       []ControlCheck `json:"control_check_fails"`
}

//ControlCheck :
type ControlCheck struct {
	CheckType      string   `json:"check_type"`
	FailedPackages int      `json:"failed_packages"`
	ListOfPackages []string `json:"failed_package_list"`
}

//ReleaseMD :
type ReleaseMD struct {
	Tag       string    `json:"release_tag"`
	CreatedAt time.Time `json:"created_at"`
	CommitID  string    `json:"commit_id"`
}

//ReleaseInsights :
type ReleaseInsights struct {
	PackageName            string        `json:"package_name"`
	CurrentVersion         string        `json:"current_version"`
	Repo                   string        `json:"repo"`
	LatestVersion          string        `json:"latest_version"`
	LatestReleaseTimestamp time.Time     `json:"latest_release_timestamp"`
	IsLatest               bool          `json:"is_latest"`
	ReleaseLag             int           `json:"release_lag"`
	MajorReleaseLag        int           `json:"major_release_lag"`
	ReleaseTimeLag         string        `json:"release_time_lag"`
	Labels                 []string      `json:"annotations"`
	ChangeInsights         CommitHistory `json:"commit_history"`
}

//RecommendedVersion :
type RecommendedVersion struct {
	PackageName        string    `json:"package_name"`
	Version            string    `json:"recommended_version"`
	ReleaseTimestamp   time.Time `json:"release_timestamp"`
	NumAuthors         int       `json:"num_uniq_authors"`
	NumReviewers       int       `json:"num_uniq_reviewers"`
	ZombieChanges      int       `json:"zombie_commits"`
	NonPeerReviewedPRs int       `json:"non_peer_reviewed_prs"`
	ChangeLabels       []string  `json:"change_annotations"`
}

//GaugeReport :
type GaugeReport struct {
	ReleaseReport struct {
		Recommendation RecommendedVersion `json:"recommendations"`
		Insights       ReleaseInsights    `json:"insights"`
	} `json:"release_report"`
}
