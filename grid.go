package spatially

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Spatially/go-geometry"

	"github.com/paulmach/go.geojson"

	"github.com/pkg/errors"
)

//
type GridPopulation struct {
	Residents          int `json:"residents"`
	Workers            int `json:"workers"`
	Housing            int `json:"housing"`
	HousingMeanValue   int `json:"housingMeanValue"`
	HousingMedianValue int `json:"housingMedianValue"`
	HousingSumValue    int `json:"housingSumValue"`
}

//
func Population(api API, locationWKT string, buffer int) (pop *GridPopulation, err error) {
	type populationResponse struct {
		Result *GridPopulation `json:"result"`
	}
	feature, err := NewFeatureFromWKT(locationWKT)
	if err != nil {
		return pop, errors.Wrap(err, "feature from wkt")
	}
	if feature.Geometry.Type != geojson.GeometryPoint {
		return pop, errors.New("location wkt must be a point")
	}
	lon, lat := feature.Geometry.Point[0], feature.Geometry.Point[1]
	url := fmt.Sprintf(SpatiallyAPI+"/grid/pop/point?lat=%v&lon=%v&radius=%v", lat, lon, buffer)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return pop, errors.Wrap(err, "prepare http request")
	}
	api.PrepareRequest(request)
	requestClient := &http.Client{}
	resp, err := requestClient.Do(request)
	if err != nil {
		return pop, errors.Wrap(err, "http get")
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return pop, errors.Wrap(err, "read response body")
	}
	if resp.StatusCode != 200 {
		return pop, api.Error(responseBody)
	}
	response := populationResponse{}
	if err = json.Unmarshal(responseBody, &response); err != nil {
		return pop, errors.Wrap(err, "parse response body json")
	}
	return response.Result, nil
}

//
func TradeAreaMarketSize(api API, tradeArea *ATA) (pop *GridPopulation, err error) {
	type tradeAreaMarketSizeRequest struct {
		FeatureCollection *geometry.FeatureCollection `json:"featureCollection"`
	}
	requestBody := tradeAreaMarketSizeRequest{tradeArea.FeatureCollection}
	j, err := json.Marshal(requestBody)
	if err != nil {
		return nil, errors.Wrap(err, "request to json")
	}
	body := bytes.NewReader(j)
	request, err := http.NewRequest("POST", SpatiallyAPI+"/grid/pop", body)
	if err != nil {
		return nil, errors.Wrap(err, "ata request")
	}
	api.PrepareRequest(request)
	requestClient := &http.Client{}
	resp, err := requestClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "request do")
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read ata request response body")
	}
	pop = &GridPopulation{}
	if err := json.Unmarshal(responseBody, pop); err != nil {
		return nil, errors.Wrap(err, "json unmarshal")
	}
	return pop, nil
}

//
type GridPopularTimes struct {
	Hours []HourlyStops `json:"result"`
}

//
type HourlyStops struct {
	Hour  int `json:"hour"`
	Stops int `json:"stops"`
}

//
func PopularTimes(api API, locationWKT string, buffer int) (pt *GridPopularTimes, err error) {
	v := &url.Values{}
	v.Add("wkt", locationWKT)
	v.Add("radius", fmt.Sprintf("%v", buffer))
	url := fmt.Sprintf(SpatiallyAPI+"/grid/stops?%v", v.Encode())
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return pt, errors.Wrap(err, "prepare http request")
	}
	api.PrepareRequest(request)
	requestClient := &http.Client{}
	resp, err := requestClient.Do(request)
	if err != nil {
		return pt, errors.Wrap(err, "http get")
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return pt, errors.Wrap(err, "read response body")
	}
	if resp.StatusCode != 200 {
		return pt, api.Error(responseBody)
	}
	pt = &GridPopularTimes{}
	if err = json.Unmarshal(responseBody, pt); err != nil {
		return pt, errors.Wrap(err, "parse response body json")
	}
	return pt, nil
}

//
type GridDistanceSensitivity struct {
	Distances Distances `json:"result"`
}

// Distances in miles
type Distances struct {
	D0  int `json:"d0"`
	D1  int `json:"d1"`
	D2  int `json:"d2"`
	D5  int `json:"d5"`
	D10 int `json:"d10"`
	D15 int `json:"d15"`
	D20 int `json:"d20"`
	D25 int `json:"d25"`
}

//
func DistanceSensitivity(api API, locationWKT string, buffer int) (ds *GridDistanceSensitivity, err error) {
	feature, err := NewFeatureFromWKT(locationWKT)
	if err != nil {
		return ds, errors.Wrap(err, "feature from wkt")
	}
	if feature.Geometry.Type != geojson.GeometryPoint {
		return ds, errors.New("location wkt must be a point")
	}
	lon, lat := feature.Geometry.Point[0], feature.Geometry.Point[1]
	request, err := http.NewRequest("GET", fmt.Sprintf(SpatiallyAPI+"/grid/distance?lat=%v&lon=%v&radius=%v", lat, lon, buffer), nil)
	if err != nil {
		return ds, errors.Wrap(err, "prepare http request")
	}
	api.PrepareRequest(request)
	requestClient := &http.Client{}
	resp, err := requestClient.Do(request)
	if err != nil {
		return ds, errors.Wrap(err, "http get")
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ds, errors.Wrap(err, "read response body")
	}
	if resp.StatusCode != 200 {
		return ds, api.Error(responseBody)
	}
	ds = &GridDistanceSensitivity{}
	if err = json.Unmarshal(responseBody, ds); err != nil {
		return ds, errors.Wrap(err, "parse response body json")
	}
	return ds, nil
}

//
type GridDemographics struct {
	TotalGeohashes  int              `json:"total"`
	WorkerTypes     []PopulationType `json:"workerTypes"`
	PopulationTypes []PopulationType `json:"populationType"`
	Demographics    map[string]int   `json:"demographics"`
}

//
type PopulationType struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

//
func Demographics(api API, locationWKT string, buffer int) (gd *GridDemographics, err error) {
	feature, err := NewFeatureFromWKT(locationWKT)
	if err != nil {
		return gd, errors.Wrap(err, "feature from wkt")
	}
	if feature.Geometry.Type != geojson.GeometryPoint {
		return gd, errors.New("location wkt must be a point")
	}
	lon, lat := feature.Geometry.Point[0], feature.Geometry.Point[1]
	request, err := http.NewRequest("GET", fmt.Sprintf(SpatiallyAPI+"/grid/highlights?lat=%v&lon=%v&radius=%v", lat, lon, buffer), nil)
	if err != nil {
		return gd, errors.Wrap(err, "prepare http request")
	}
	api.PrepareRequest(request)
	requestClient := &http.Client{}
	resp, err := requestClient.Do(request)
	if err != nil {
		return gd, errors.Wrap(err, "http get")
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return gd, errors.Wrap(err, "read response body")
	}
	if resp.StatusCode != 200 {
		return gd, api.Error(responseBody)
	}
	type demographicsResponse struct {
		Result *GridDemographics `json:"result"`
	}
	response := &demographicsResponse{}
	if err = json.Unmarshal(responseBody, response); err != nil {
		return gd, errors.Wrap(err, "parse response body json")
	}
	return response.Result, nil
}
