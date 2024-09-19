package Type

type QueryData struct {
	StoreName    string
	StoreAddress string
	StoreRating  float64
	StoreMapLink string
}

func (q *QueryData) Init(storeName string, storeAddress string, storeRating float64, storeMapLink string) {
	q.StoreName = storeName
	q.StoreAddress = storeAddress
	q.StoreRating = storeRating
	q.StoreMapLink = storeMapLink
}
