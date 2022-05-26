package common

//GithubRateLimitErr :
type GithubRateLimitErr struct{}

func (m *GithubRateLimitErr) Error() string {
	return "API rate limit reached"
}
