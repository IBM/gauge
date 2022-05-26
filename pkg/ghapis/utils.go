package ghapis

import "strings"

func parseRepositoryOwner(giturl string) string {
	var owner string
	gParts := strings.Split(giturl, "/")
	if len(gParts) >= 4 {
		owner = gParts[3]
	}
	return owner
}

func parseRepositoryName(giturl string) string {
	var repo string
	repoURL := strings.TrimRight(giturl, "/")
	urlParts := strings.Split(repoURL, "/")
	repo = urlParts[len(urlParts)-1]
	return repo
}

func contains(cache map[int]struct{}, key int) bool {
	if _, ok := cache[key]; ok {
		return true
	}
	return false
}
