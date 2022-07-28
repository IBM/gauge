## Gauge

Measure release insights and get recommendations for open-source dependencies.

## About Gauge

For OSS ecosystem, there are established practices for disclosing, discovering and remediating vulnerabilities in the code. Although, in the wake of recent cybersecurity incidents it is becoming important to understand and assess risks associated with developers and their contribution practices. Project “gauge” aims to provides risk assessment for release engineering. For instance, when you upgrade your OSS dependency, it measures risk from every commit that went into the new release, developers that contributed those changes, code review practices observed and types of changes that went into release (performance fix, security fix, but fix, etc.). Core motivation behind project gauge is to bring visibility and auditing into OSS releases.

## Requirements

To be able to run this project successfully, the following needs to be configured:

1. First, Golang must be installed. Downloaders can be installed from the Go website located here: https://go.dev/doc/install

2. Set necessary environment variables.
- `GITHUB_API_KEY` needs to be set as an environment variable with a valid PAT token, with the "public_repo" scope defined. For details on how to set this up, documentation is located here: https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token. 

> We have also setup a cache to avoid making repeated API calls to GitHub. This should help with better performance and get around hitting the API rate limit. Please open an issue to get access to the cache API server. Soon we plan to make it publicly accessible.

- `RELEASE_LIB_SERVER` (optional): Set the caching API server URL 


3. Next, the project can be compiled. To build the project, compile the code with the following command to create the executable `gauge`: 
```
make
```

## Config file

Copy the .gauge.yaml file and update the control check parameters/thresholds. 

```yaml
## gauge control file

runtime-configs:
  # releaselib service
  releaselib-service: "http://127.0.0.1:9950"

  # weather api key
  # weather apis are used to resolve location information
  # from github to specific country names
  weather-api-key: 

  # github api key
  # github api key is use to avoid rate limit
  github-api-key: 

release-control:
  enable: true
  # maximum release lag for dependencies in terms of
  # their versions 
  max-release-lag: 3

  # maximum release lag for dependencies in terms of
  # time duration (in days)
  max-release-lag-duration: 180

  # Ensure every code change (pull request) has been reviewed
  # by atleast one reviewers (who is different from the author)
  peer-review-enforced: true

  # Zombie changes are the ones that are commited to `main` 
  # branch directly without formal pull request
  # Control to block such code changes
  zombie-commit-enforced: true

## export control check verifies developers/contributors location against
## known export control restricted coutries
export-controls:
  # flag to enable/disable export control check 
  enable: true

  # list of countries to check against
  restricted-countries: [country-name-1]

  # countribution threshold
  contribution-threshold: 10

  # regulated compliance control country list
  taa-list: []

  ofac-list: []
```

## Sample Run

Next, its time to try it out! 

There are two operational modes for running gauge today. 

1. Package: Evaluate health of individual package/repository

2. SBOM : Evaluate health of all OSS packages from the SBOM

### Package Mode

Say, you are using python package `flask` with current version `2.1.1` and you want to evaluate next version before you upgrade to it. 
You can run following query against `gauge` to get those insights:

```
./gauge package -p flask -e python -t 2.1.1 -r https://github.com/pallets/flask

complete log file is available at: /tmp/gauge-075403876
********************************************************************************
Gauge Report for package `flask`
********************************************************************************
Release Measures:
	Current version: 2.1.1
	Latest version: 2.1.2
	Release lag (versions): 1
	Release lag (days): 28 days
--------------------------------------------------------------------------------
		Recommended update
		 Version - 2.1.2
		 Release Time - 2022-04-28 17:48:24 +0000 UTC
		 Num of unique contributors - 7
		 Num of unique reviewers - 0
		 Non peer reviewed changes - 0
		 Num of zombie commits - 15
		 Change annotations - ['docs','typing','testing']
--------------------------------------------------------------------------------                 
```

In this mode, we have limited feature to discover (if missing) github repsitory path from package-names. 

### SBOM Mode

> This will be shortly available

In this mode, `gauge` can accept CycloneDX/SPDX formatted SBOM as input and provide evaluation/recommendations for every OSS dependency from SBOM. 


## WIP

We envision `gauge` to cover the OSS universe that also includes different modalities of packaging/distributing OSS components: 
 1. container images
 2. Deployment YAMLs (e.g. CRDs, policies, tekton tasks)
 3. App bundles (e.g. k8s operators)

We also have an opportunity to enrich our release insights and recommendations with sophistcated ML techniques.
