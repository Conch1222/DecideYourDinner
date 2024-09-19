package Type

type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type GeoLocation struct {
	Location Location `json:"location"`
	Accuracy float64  `json:"accuracy"`
}

type NearBy struct {
	NearByResults []NearByResult `json:"results"`
	Status        string         `json:"status"`

	Option string
	Weight float64
}

type NearByResult struct {
	BusinessStatus  string       `json:"business_status"`
	Name            string       `json:"name"`
	PlaceId         string       `json:"place_id"`
	PriceLevel      int          `json:"price_level"`
	Rating          float64      `json:"rating"`
	Reference       string       `json:"reference"`
	UserRatingTotal int          `json:"user_ratings_total"`
	Vicinity        string       `json:"vicinity"`
	Geometry        GeoLocation  `json:"geometry"`
	PlusCode        PlusCode     `json:"plus_code"`
	OpeningHours    OpeningHours `json:"opening_hours"`

	RankingScore float64
}

type Geometry struct {
	Location Location `json:"location"`
	Viewport Viewport `json:"viewport"`
}

type Viewport struct {
	Northeast Location `json:"northeast"`
	Southwest Location `json:"southwest"`
}

type PlusCode struct {
	CompoundCode string `json:"compound_code"`
	GlobalCode   string `json:"global_code"`
}

type OpeningHours struct {
	OpenNow bool `json:"open_now"`
}
