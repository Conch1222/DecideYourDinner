package File

type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type GeoLocation struct {
	Location Location `json:"location"`
	Accuracy float64  `json:"accuracy"`
}
