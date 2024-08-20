package main

import (
	"GoWeb/Error"
	"GoWeb/File"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type GMapClient struct {
	Client *http.Client
	ApiKey string
}

var Client *GMapClient
var once sync.Once

func getClient() *GMapClient {
	// singleton
	once.Do(func() {
		client, _ := initClient()
		Client = client
	})
	return Client
}

func initClient() (*GMapClient, error) {
	key, err := ReadApiKey()
	if err != nil {
		panic(err)
	}

	Client = new(GMapClient)
	Client.Client = &http.Client{Timeout: time.Second * 10}
	Client.ApiKey = key
	return Client, nil
}

func ReadApiKey() (string, error) {
	file, err := os.Open("File/Key.txt")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	return "", errors.New(Error.InvalidApiKey)
}

func (client *GMapClient) getUserLocation(w http.ResponseWriter) (File.GeoLocation, error) {
	// POST
	dst := "https://www.googleapis.com/geolocation/v1/geolocate?key=" + client.ApiKey
	form := url.Values{}

	req, err := http.NewRequest("POST", dst, strings.NewReader(form.Encode()))
	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Client.Do(req)
	data, _ := io.ReadAll(res.Body)
	fmt.Fprintln(w, string(data))

	var resultLocation File.GeoLocation
	if err := json.Unmarshal(data, &resultLocation); err == nil {
		return resultLocation, nil
	} else {
		return resultLocation, errors.New(Error.CannotGetUserLocation)
	}
}

func (client *GMapClient) getUserNearBy(w http.ResponseWriter, loc File.GeoLocation, option string) {
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
	q.Add("radius", strconv.Itoa(100))
	q.Add("key", client.ApiKey)
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())

	res, err := client.Client.Do(req)
	data, _ := io.ReadAll(res.Body)
	fmt.Fprintln(w, string(data))
}

func convertGeoLocationToString(location File.GeoLocation) string {
	lat := strconv.FormatFloat(location.Location.Latitude, 'f', -1, 64)
	lng := strconv.FormatFloat(location.Location.Longitude, 'f', -1, 64)
	return lat + "," + lng
}
