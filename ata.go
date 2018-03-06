package spatially

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Spatially/go-geometry"
	"github.com/pkg/errors"
)

//
type ATA struct {
	*geometry.FeatureCollection
}

//
type ATALocationType int

//
const (
	Home ATALocationType = iota
	Work
	HomeAndWork
)

//
type ATATimeOfDay int

//
const (
	AllDay ATATimeOfDay = iota
	Morning
	MidDay
	Evening
	Night
)

func (a ATATimeOfDay) String() string {
	switch a {
	case AllDay:
		return "AllDay"
	case Morning:
		return "Morning"
	case MidDay:
		return "MidDay"
	case Evening:
		return "Evening"
	case Night:
		return "Night"
	default:
		return "AllDay"
	}
}

//
type ATAOptions struct {
	LocationType ATALocationType
	TimeOfDay    ATATimeOfDay
	GeoFence     bool
}

type ataRequest struct {
	PointWKT     string   `json:"pointWKT"`
	AreaType     string   `json:"areaType"`
	Buffer       int      `json:"buffer"`
	Distance     int      `json:"distance"`
	TimeOfDay    string   `json:"timeOfDay"`
	LocationType []string `json:"locationType"`
	GeoFence     bool     `json:"geoFence"`
}

type ataResponse struct {
	Buffer            int                         `json:"int"`
	FeatureCollection *geometry.FeatureCollection `json:"featureCollection"`
	Version           int                         `json:"version"`
	Message           string                      `json:"message"`
}

//
func NewATA(api API, locationWKT string, options *ATAOptions) (ata *ATA, err error) {
	requestBody := &ataRequest{
		PointWKT: locationWKT,
		AreaType: "ATA",
		Buffer:   100,
		Distance: 0,
	}
	if options != nil {
		requestBody.TimeOfDay = options.TimeOfDay.String()
		switch options.LocationType {
		case Home:
			requestBody.LocationType = []string{"Home"}
		case Work:
			requestBody.LocationType = []string{"Work"}
		case HomeAndWork:
			requestBody.LocationType = []string{"Home", "Work"}
		}
		requestBody.GeoFence = options.GeoFence
	} else {
		requestBody.TimeOfDay = AllDay.String()
		requestBody.LocationType = []string{"Home"}
	}
	j, err := json.Marshal(requestBody)
	if err != nil {
		return nil, errors.Wrap(err, "request to json")
	}
	body := bytes.NewReader(j)
	request, err := http.NewRequest("POST", SpatiallyAPI+"/ads/science/ata", body)
	if err != nil {
		return nil, errors.Wrap(err, "ata request")
	}
	api.PrepareRequest(request)
	requestClient := &http.Client{}
	resp, err := requestClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "ata request do")
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read ata request response body")
	}
	var response ataResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, errors.Wrap(err, "json unmarshal ata request")
	}
	if response.FeatureCollection == nil {
		return nil, errors.New(response.Message)
	}
	ata = &ATA{response.FeatureCollection}
	return
}
