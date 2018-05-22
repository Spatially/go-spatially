// Package spatially provides a set of APIs to communicate with Spatially DB API to create layers
// and features in your database. It supports operation with GeoJSON and Well-Known-Text feature types
package spatially

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

//
type API interface {
	PrepareRequest(req *http.Request)
	Error(response []byte) (err error)
}

// SpatialConstraint is an object used to describe and boundary and intersection type from which
// to query features with
type SpatialConstraint struct {
	WKT    string                `json:"wkt"`
	Radius float64               `json:"radius"`
	Type   SpatialConstraintType `json:"type"`
}

// SpatialConstraintType is the type of spatial intersection to do on features
type SpatialConstraintType int

const (
	// SpatialConstraintIntersect is a SpatialContraintType that only selects features that intersect
	// with the given boundary
	SpatialConstraintIntersect SpatialConstraintType = iota
)

// SpatiallyAPI - Spatially's API URL
const SpatiallyAPI = "https://api.spatially.com"

type api struct {
	Token string
}

type gatewayRequest struct {
	Code string `json:"code"`
	Key  string `json:"key"`
}

type gatewayResponse struct {
	Token string `json:"token"`
}

// New created a new instance of the Spatially API. The parameters are the api code & key
// provided by Spatially. It generates a token with the API.
func NewAPI(apiCode, apiKey string) (API, error) {
	request := gatewayRequest{
		Code: apiCode,
		Key:  apiKey,
	}
	j, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb json marshal gateway request")
	}
	body := bytes.NewReader(j)
	resp, err := http.Post(SpatiallyAPI+"/gateway/client", "application/json", body)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb gateway request")
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "spatialdb read gateway request response body")
	}
	var response gatewayResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, errors.Wrap(err, "spatialdb json unmarshal gateway request")
	}
	if len(response.Token) == 0 {
		return nil, errors.New("SpatialDB was not able to generate a valid token")
	}
	return &api{Token: response.Token}, nil
}

func (s *api) PrepareRequest(request *http.Request) {
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+s.Token)
}

type requestError struct {
	Message string `json:"message"`
}

func (s *api) Error(responseBody []byte) error {
	var err requestError
	if err := json.Unmarshal(responseBody, &err); err != nil {
		return errors.Wrap(errors.New(string(responseBody)), "unable to parse error body")
	}
	return errors.New(err.Message)
}
