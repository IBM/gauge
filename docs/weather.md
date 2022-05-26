# Weather API Algorithm 
## Overview 
The algorithm used for country resolution is located at [../pkg/weather/client.go](../pkg/weather/client.go)

The API used in this project, Weather Company's LocationSearch API is documented here: https://weather.com/swagger-docs/ui/sun/v3/sunV3LocationSearch.json 

## Details

The function that gets called from the main logic of the code is GetLocation() 

The first thing it does is remove special characters from the string. If there are certain special characters, then the Weather API will return 400 error Bad Request. 

Next, we detemine if the location is type "not available", which means we know there is nothing in the location field of that Github user's profile. 

The LocationSearch API takes in a few parameters – “query” and ”LocationType”. “query” is the location text we get from the Github user – “LocationType” is the type of location of the query, options include “City” “State” or “Country” and others. The API returns various location information broken down into types, but for our purposes we take the “Country” field returned and use this for the API.  

A challenge we have here is that the “LocationType” is not known when we make our request. So what we did was for each location, we make three requests, one for LocationType = City, State, and Country. We then make our best guess from the information returned from these three requests. 

Next we make make a series of if statements to best guess the country: 
1. If all requests return nothing, return nothing.
2. If only one request is not null, return the single not null result.
3. If both city and state type requests are not null, return city results if the two results are not the same, and return the shared result if they are. 
4. If both city and country type results are not null, return country if the two results are not the same, and return the shared result if they are. 
5. If both state and country type results are not null, return country if the two results are not the same, and return the shared result if they are. 
6.  If all results are the same, return the shared result.
7. If city and state are the same result, return city result.
8. If country result the same as city or state, return city result. 
9. If all results are different and not empty, then return the country if the location string contains a comma, and the city if not.  

The country then gets returned, and in the main logic the results get cached to reduce the load to the api. 
