package pkgmgr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//following repo solution code is from scorecard project

type npmSearchResults struct {
	Objects []struct {
		Package struct {
			Links struct {
				Repository string `json:"repository"`
			} `json:"links"`
		} `json:"package"`
	} `json:"objects"`
}

type pypiSearchResults struct {
	Info struct {
		ProjectUrls struct {
			Source string `json:"Source Code"`
		} `json:"project_urls"`
	} `json:"info"`
}

// FetchGitRepositoryFromNPM : Gets the GitHub repository URL for the npm package.
// nolint: noctx
func FetchGitRepositoryFromNPM(packageName string) (string, error) {
	npmSearchURL := "https://registry.npmjs.org/-/v1/search?text=%s&size=1"
	const timeout = 10
	client := &http.Client{
		Timeout: timeout * time.Second,
	}
	resp, err := client.Get(fmt.Sprintf(npmSearchURL, packageName))
	if err != nil {
		return "", fmt.Errorf("failed to get npm package json: %v", err)
	}

	defer resp.Body.Close()
	v := &npmSearchResults{}
	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return "", fmt.Errorf("failed to parse npm package json: %v", err)
	}
	if len(v.Objects) == 0 {
		return "",
			fmt.Errorf("could not find source repo for npm package: %s", packageName)
	}
	return v.Objects[0].Package.Links.Repository, nil
}

// FetchGitRepositoryFromPYPI :Gets the GitHub repository URL for the pypi package.
// nolint: noctx
func FetchGitRepositoryFromPYPI(packageName string) (string, error) {
	pypiSearchURL := "https://pypi.org/pypi/%s/json"
	const timeout = 10
	client := &http.Client{
		Timeout: timeout * time.Second,
	}
	resp, err := client.Get(fmt.Sprintf(pypiSearchURL, packageName))
	if err != nil {
		return "", fmt.Errorf("failed to get pypi package json: %v", err)
	}

	defer resp.Body.Close()
	v := &pypiSearchResults{}
	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return "", fmt.Errorf("failed to parse pypi package json: %v", err)
	}
	if v.Info.ProjectUrls.Source == "" {
		return "", fmt.Errorf("could not find source repo for pypi package: %s", packageName)
	}
	return v.Info.ProjectUrls.Source, nil
}
