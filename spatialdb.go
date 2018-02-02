package spatialdb

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

//
type SpatialDB interface {
	// Layers
	GetLayers() (layers Layers, err error)
	GetLayer(id string) (layer *Layer, err error)
	CreateLayer(name string) (layer *Layer, err error)
	DeleteLayer(id string) (err error)

	// Features
	GetFeatures(layer string, spatialConstraint *SpatialConstraint) (features Features, err error)
	GetFeature(id string) (feature *Feature, err error)
	CreateFeature(layer string, feature *Feature) (createdFeature *Feature, err error)
	UpdateFeature(id string, properties map[string]interface{}) (feature *Feature, err error)
	DeleteFeature(id string) (err error)
}

//
type SpatialConstraint struct {
	WKT    string                `json:"wkt"`
	Radius float64               `json:"radius"`
	Type   SpatialConstraintType `json:"type"`
}

//
type SpatialConstraintType int

//
const (
	SpatialConstraintIntersect SpatialConstraintType = iota
	SpatialConstraintWithin                          // TODO(josebalius): Support this
)

//
const SpatiallyAPI = "http://localhost:8000"

type spatialDB struct {
	Token string
}

type gatewayRequest struct {
	Code string `json:"code"`
	Key  string `json:"key"`
}

type gatewayResponse struct {
	Token string `json:"token"`
}

// TODO(josebalius): Name ALL errors?
// TODO(josebalius): httptest mock server
func New(apiCode, apiKey string) (SpatialDB, error) {
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
	return spatialDB{Token: response.Token}, nil
}

func (s spatialDB) prepareRequest(request *http.Request) {
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+s.Token)
}

type requestError struct {
	Message string `json:"message"`
}

func (s spatialDB) requestError(responseBody []byte) error {
	var err requestError
	if err := json.Unmarshal(responseBody, &err); err != nil {
		return err
	}
	return errors.New(err.Message)
}
