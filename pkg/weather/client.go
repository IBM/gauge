package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/IBM/gauge/pkg/common"
)

type WeatherResp struct {
	Location LocationStruct `json:"location"`
}

type LocationStruct struct {
	Address []string `json:"address"`
	City    []string `json:"city"`
	Country []string `json:"country"`
}

func SearchCity(location string) (string, error) {
	locationType := "city"
	apikey := os.Getenv(common.WEATHER_API_KEY)
	url := fmt.Sprintf("https://api.weather.com/v3/location/search?query=%s&locationType=%s&language=en-US&format=json&apiKey=%s", url.PathEscape(location), locationType, apikey)
	retries := 3
	var respobject WeatherResp
	var statuscode int
	for try := 1; try <= retries; try++ {
		resp, err := http.Get(url)
		if err != nil {
			if try == retries {
				return "", errors.New("Error making Weather API call")
			}
			time.Sleep(3 * time.Second)
			continue
		}
		statuscode = resp.StatusCode
		if statuscode == 404 {
			// location not found
			return "", nil
		} else if statuscode == 503 {
			return "", fmt.Errorf("503 network error")
		} else if statuscode == 400 {
			// Network Error, Bad Request, usually invalid characters
			return "", nil
		} else if statuscode != 200 {
			// error making api call
			return "", fmt.Errorf("%d error", statuscode)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &respobject); err != nil {
			return location, fmt.Errorf("error unmarshaling response err: %v", err)
		}
		break
	}
	country := respobject.Location.Country[0]
	return country, nil
}

func SearchState(location string) (string, error) {
	locationType := "state"
	apikey := os.Getenv(common.WEATHER_API_KEY)
	url := fmt.Sprintf("https://api.weather.com/v3/location/search?query=%s&locationType=%s&language=en-US&format=json&apiKey=%s", url.PathEscape(location), locationType, apikey)
	retries := 3
	var respobject WeatherResp
	var statuscode int
	for try := 1; try <= retries; try++ {
		resp, err := http.Get(url)
		if err != nil {
			if try == retries {
				return "", errors.New("Error making Weather API call")
			}
			time.Sleep(3 * time.Second)
			continue
		}
		statuscode = resp.StatusCode
		if statuscode == 404 {
			// location not found
			return "", nil
		} else if statuscode == 401 {
			if try == retries {
				return "", errors.New("401 Error: reached Weather API rate limit error")
			}
			// rate limit error, try again
			time.Sleep(3 * time.Second)
			continue
		} else if statuscode == 503 {
			return "", fmt.Errorf("503 network error")
		} else if statuscode == 400 {
			// Network Error, Bad Request, usually invalid characters
			return "", nil
		} else if statuscode != 200 {
			// error making api call
			return "", fmt.Errorf("%d error", statuscode)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &respobject); err != nil {
			return location, fmt.Errorf("error unmarshaling response err: %v", err)
		}
		break
	}
	country := respobject.Location.Country[0]
	return country, nil
}

func SearchCountry(location string) (string, error) {
	locationType := "country"
	apikey := os.Getenv(common.WEATHER_API_KEY)
	url := fmt.Sprintf("https://api.weather.com/v3/location/search?query=%s&locationType=%s&language=en-US&format=json&apiKey=%s", url.PathEscape(location), locationType, apikey)
	retries := 3
	var respobject WeatherResp
	var statuscode int
	for try := 1; try <= retries; try++ {
		resp, err := http.Get(url)
		if err != nil {
			if try == retries {
				return "", errors.New("Error making Weather API call")
			}
			time.Sleep(5 * time.Second)
			continue
		}
		statuscode = resp.StatusCode
		if statuscode == 404 {
			// location not found
			return "", nil
		} else if statuscode == 503 {
			return "", fmt.Errorf("503 network error")
		} else if statuscode == 400 {
			// Network Error, Bad Request, usually invalid characters
			return "", nil
		} else if statuscode != 200 {
			// error making api call
			return "", fmt.Errorf("%d error", statuscode)
		}

		body, _ := ioutil.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &respobject); err != nil {
			return location, fmt.Errorf("error unmarshaling response err: %v", err)
		}
		break
	}
	country := respobject.Location.Country[0]
	return country, nil
}

func RemoveCharacters(input string, characters string) string {
	filter := func(r rune) rune {
		if strings.IndexRune(characters, r) < 0 {
			return r
		}
		return -1
	}
	return strings.Map(filter, input)

}

func GetLocation(location string) (string, error) {
	location = RemoveCharacters(location, "<>;&|\"{}\\^%")
	splitlocation := strings.Split(location, ",")

	if location == "navail" || location == "not available" || location == "" {
		return "not available", nil
	} else if location == "" {
		return "not available", nil
	}

	cityres, err := SearchCity(location)
	if err != nil {
		return cityres, err
	}
	// time.Sleep(600 * time.Millisecond)
	stateres, err := SearchState(location)
	if err != nil {
		return stateres, err
	}
	// time.Sleep(600 * time.Millisecond)
	countryres, err := SearchCountry(location)
	if err != nil {
		return countryres, err
	}
	// time.Sleep(600 * time.Millisecond)

	if cityres == "" && stateres == "" && countryres == "" {
		return "", nil
	} else if cityres == "" && stateres == "" && countryres != "" {
		return countryres, nil
	} else if cityres == "" && stateres != "" && countryres == "" {
		return stateres, nil
	} else if cityres != "" && stateres == "" && countryres == "" {
		return cityres, nil
	} else if cityres != "" && stateres != "" && countryres == "" {
		if cityres == stateres {
			return stateres, nil
		} else {
			return cityres, nil
		}
	} else if cityres != "" && stateres == "" && countryres != "" {
		if cityres == countryres {
			return cityres, nil
		} else {
			return countryres, nil
		}
	} else if cityres == "" && stateres != "" && countryres != "" {
		if stateres == countryres {
			return stateres, nil
		} else {
			return countryres, nil
		}
	}
	// at this point none of the results are empty
	if cityres == stateres && stateres == countryres {
		return countryres, nil
	} else if cityres == stateres {
		return cityres, nil
	} else if cityres == countryres {
		return countryres, nil
	} else if stateres == countryres {
		return countryres, nil
	}
	//at this point all results are not empty and different

	// describe the logic here
	if len(splitlocation) == 1 {
		return countryres, nil
	} else {
		return cityres, nil
	}
}
