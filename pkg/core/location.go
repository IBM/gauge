package core

import (
	"github.com/IBM/gauge/pkg/common"
)

func resolveLocation(contributors []common.ContributorMD) common.ErrorReport {

	errorReport := common.ErrorReport{}
	var resolvedCountry string
	var err error
	resolverReady := false

	_, resolverReady = os.LookupEnv(common.WEATHER_API_KEY)
	for idx, cloc := range contributors {
		cleanLocation := weather.RemoveCharacters(cloc.Location, "<>;&|\"{}")
		if cleanLocation == "not available" || cleanLocation == "navail" || cleanLocation == "" {
			errorReport.Reasons = append(errorReport.Reasons, fmt.Sprintf("no location data in profile for contributor: `%s`", cloc.LoginID))
			errorReport.NumberOfError++
			continue
		}
		if !resolverReady {
			contributors[idx].ResolvedCountry = cleanLocation
			continue
		}
		if cloc.ResolvedCountry == "" {
			resolvedCountry, err = releaselib.GetLocationMeta(cleanLocation)
			if err != nil {
				resolvedCountry, err = weather.GetLocation(cleanLocation)
				if err != nil {
					errorReport.Reasons = append(errorReport.Reasons, fmt.Sprintf("error resolving location `%s`: %s ", cloc.Location, err))
					errorReport.NumberOfError++
				}
				if resolvedCountry == "" {
					resolvedCountry = "not resolved"
				}
			}
			releaselib.StoreLocationMeta(cleanLocation, resolvedCountry, cloc.LoginID)
			contributors[idx].ResolvedCountry = resolvedCountry
		} else {
			resolvedCountry = cloc.ResolvedCountry
		}
		if resolvedCountry == "not resolved" {
			// fmt.Printf("location data not resolved for contributor: \"%s\", URL: \"%s\", location: \"%s\"\n", cloc.Name, cloc.URL, cloc.Location)
			errorReport.Reasons = append(errorReport.Reasons, fmt.Sprintf("location data not resolved for contributor: `%s`, location: `%s`", cloc.Name, cloc.Location))
			errorReport.NumberOfError++
			// continue
		}
	}
	return errorReport
}
