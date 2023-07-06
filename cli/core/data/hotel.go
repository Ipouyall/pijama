package data

type HotelRoom struct {
	ID         int    `json:"id"`
	HotelName  string `json:"hotel_name"`
	HotelClass string `json:"hotel_class"`
	Cost       int    `json:"cost"`
	City       string `json:"city"`
	Address    string `json:"address"`
}
