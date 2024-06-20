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
	file, err := os.Open("File/key.txt")
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
