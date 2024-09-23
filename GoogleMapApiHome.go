package main

import (
	"GoWeb/Error"
	"GoWeb/Type"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type GMapClient struct {
	Client *http.Client
	ApiKey string
}

var GClient *GMapClient
var onceGClient sync.Once

func getClient() (*GMapClient, error) {
	// singleton
	error := error(nil)
	onceGClient.Do(func() {
		fmt.Println("once Init Cli")
		client, err := initClient()
		if err != nil {
			error = err
		}
		GClient = client
	})
	if error != nil {
		return nil, error
	}
	return GClient, nil
}

func initClient() (*GMapClient, error) {
	key, err := ReadKey("File/Key.txt", Error.InvalidApiKey)
	if err != nil {
		fmt.Println(Error.InvalidApiKey)
		return nil, err
	}

	if GClient == nil {
		fmt.Println("init client")
		GClient = new(GMapClient)
		GClient.Client = &http.Client{Timeout: time.Second * 10}
		GClient.ApiKey = key
		return GClient, nil
	}

	return GClient, nil
}

func (client *GMapClient) getUserLocation(w http.ResponseWriter) (*Type.GeoLocation, error) {
	// POST
	dst := "https://www.googleapis.com/geolocation/v1/geolocate?key=" + client.ApiKey
	form := url.Values{}

	req, err := http.NewRequest("POST", dst, strings.NewReader(form.Encode()))
	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Client.Do(req)
	data, _ := io.ReadAll(res.Body)

	var resultLocation Type.GeoLocation
	if err := json.Unmarshal(data, &resultLocation); err == nil {
		return &resultLocation, nil
	} else {
		return nil, errors.New(Error.OutputError_CannotGetUserLocation)
	}
}

func (client *GMapClient) getUserNearBy(w http.ResponseWriter, loc Type.GeoLocation, option string) (*Type.NearBy, error) {
	// GET
	dst := "https://maps.googleapis.com/maps/api/place/nearbysearch/json"

	req, err := http.NewRequest("GET", dst, nil)
	if err != nil {
		fmt.Println(err)
	}

	q := req.URL.Query()
	q.Add("location", convertGeoLocationToString(loc))
	q.Add("language", "zh-TW")
	q.Add("keyword", option)
	q.Add("radius", strconv.Itoa(int(loc.Accuracy+100)))
	q.Add("key", client.ApiKey)
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())

	res, err := client.Client.Do(req)
	data, _ := io.ReadAll(res.Body)

	var resultNearBy Type.NearBy
	if err := json.Unmarshal(data, &resultNearBy); err == nil {
		return &resultNearBy, nil
	} else {
		return nil, errors.New(Error.OutputError_CannotGetNearBy)
	}
}

func convertGeoLocationToString(location Type.GeoLocation) string {
	lat := strconv.FormatFloat(location.Location.Latitude, 'f', -1, 64)
	lng := strconv.FormatFloat(location.Location.Longitude, 'f', -1, 64)
	return lat + "," + lng
}
