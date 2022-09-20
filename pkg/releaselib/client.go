package releaselib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/IBM/gauge/pkg/common"
)

const (
	getPkgURI       = "/v1/packages"
	contributorsURI = "/v1/contributors"
	locationURI     = "/v1/locations"
	repositoryURI   = "/v1/repository"
	summaryURI      = "/v1/summary"
)

type ResolvedLocationCacheStruct struct {
	RawLocation     string `json:"raw_location"`
	ResolvedCountry string `json:"resolved_country"`
	ContributorID   string `json:"contributor_login_id"`
}

//StorePackageMeta :
func StorePackageMeta(pkg *common.PackageMD) error {
	relServer, cacheReady := os.LookupEnv(common.RELEASE_LIB_SERVER)
	if !cacheReady {
		return fmt.Errorf("cache server not setup")
	}
	url := fmt.Sprintf("%s%s", relServer, getPkgURI)
	payload, _ := json.MarshalIndent(pkg, "", "    ")
	retCode, _, err := common.MakePOSTAPICall(url, payload)
	if err != nil && retCode != 201 {
		return fmt.Errorf("un-expected result: retcode: %d err: %v", retCode, err)
	}
	return nil
}

//StoreRepositoryMeta :
func StoreRepositoryMeta(repoMD *common.RepositoryMD) error {
	relServer, cacheReady := os.LookupEnv(common.RELEASE_LIB_SERVER)
	if !cacheReady {
		return fmt.Errorf("cache server not setup")
	}
	url := fmt.Sprintf("%s%s", relServer, repositoryURI)
	payload, _ := json.MarshalIndent(repoMD, "", "    ")
	retCode, _, err := common.MakePOSTAPICall(url, payload)
	if err != nil && retCode != 201 {
		return fmt.Errorf("un-expected result: retcode: %d err: %v", retCode, err)
	}
	return nil
}

//GetPackageMeta :
func GetPackageMeta(pkgmd *common.PackageRepoMD) error {
	relServer, cacheReady := os.LookupEnv(common.RELEASE_LIB_SERVER)
	if !cacheReady {
		return fmt.Errorf("cache server not setup")
	}
	payload, _ := json.MarshalIndent(pkgmd, "", "    ")
	retCode, respBody, err := common.MakeGetAPICall(relServer, getPkgURI, payload)
	if err != nil && retCode != 200 {
		return fmt.Errorf("un-expected result: retcode: %d err: %v", retCode, err)
	}

	if err = json.Unmarshal(respBody, &pkgmd); err != nil {
		return fmt.Errorf("error unmarshaling response err: %v", err)
	}
	return nil
}

//GetContributorMeta :
func GetContributorMeta(pkgMeta *common.PackageRepoMD) error {
	relServer, cacheReady := os.LookupEnv(common.RELEASE_LIB_SERVER)
	if !cacheReady {
		return fmt.Errorf("cache server not setup")
	}

	type payload struct {
		RepoURL string `json:"repo_url"`
	}

	p := payload{}
	p.RepoURL = pkgMeta.RepoURL

	pBuf, _ := json.MarshalIndent(p, "", "    ")
	retCode, respBody, err := common.MakeGetAPICall(relServer, contributorsURI, pBuf)
	if err != nil && retCode != 200 {
		return fmt.Errorf("un-expected result: retcode: %d err: %v", retCode, err)
	}

	if err = json.Unmarshal(respBody, pkgMeta); err != nil {
		return fmt.Errorf("error unmarshaling response err: %v", err)
	}
	return nil
}

//StoreContributors :
func StoreContributors(payload *common.PackageRepoMD) error {
	if len(payload.Contributors) == 0 {
		return nil
	}
	relServer, cacheReady := os.LookupEnv(common.RELEASE_LIB_SERVER)
	if !cacheReady {
		return fmt.Errorf("cache server not setup")
	}
	url := fmt.Sprintf("%s%s", relServer, contributorsURI)
	pBuf, _ := json.MarshalIndent(payload, "", "    ")
	// fmt.Println(string(pBuf))
	retCode, _, err := common.MakePOSTAPICall(url, pBuf)
	if err != nil || retCode != http.StatusCreated {
		fmt.Println(err, retCode)
		return fmt.Errorf("error storing contributor metadata to releaselib")
	}
	return nil
}

//StoreLocation :
func StoreLocationMeta(rawLocation, resolvedCountry, contributorID string) error {
	relServer, cacheReady := os.LookupEnv(common.RELEASE_LIB_SERVER)
	if !cacheReady {
		return fmt.Errorf("cache server not setup")
	}

	dataobject := ResolvedLocationCacheStruct{
		RawLocation:     rawLocation,
		ResolvedCountry: resolvedCountry,
		ContributorID:   contributorID}
	payload, _ := json.Marshal(dataobject)
	// fmt.Printf("TEST: Storing location %+v\n", dataobject)
	url := fmt.Sprintf("%s%s", relServer, locationURI)
	retCode, _, err := common.MakePOSTAPICall(url, payload)
	if err != nil && retCode != 201 {
		return fmt.Errorf("un-expected result: retcode: %d err: %v", retCode, err)
	}
	return nil
}

//GetCacheState :
func GetCacheState(pkgName, ecosystem, repo string, result *common.PackageRepoMD) error {
	relServer, cacheReady := os.LookupEnv(common.RELEASE_LIB_SERVER)
	if !cacheReady {
		return fmt.Errorf("cache server not setup")
	}

	type payload struct {
		PackageName string `json:"package_name"`
		Ecosystem   string `json:"ecosystem"`
		RepoURL     string `json:"repo_url"`
	}

	p := payload{}
	p.RepoURL = repo
	p.Ecosystem = ecosystem
	p.PackageName = pkgName

	pBuf, _ := json.MarshalIndent(p, "", "    ")
	retCode, respBody, err := common.MakeGetAPICall(relServer, summaryURI, pBuf)
	if err != nil && retCode != 200 {
		return fmt.Errorf("un-expected result: retcode: %d err: %v", retCode, err)
	}
	if err = json.Unmarshal(respBody, result); err != nil {
		return fmt.Errorf("error unmarshaling response err: %v", err)
	}
	return nil
}
