package core

import (
	"saaj/core/data"
)

type Core interface {
	Authenticate(username, password string) (bool, string)

	GetPackage() []data.Package
	RequestPackage(packID int) data.Requirement

	SubmitDocument(docID int, name, content string) bool

	GetHotels() []data.HotelRoom
	ReserveHotel(hotelID int) bool

	RequestVisa() bool
	SubmitVisa(visaID int) bool

	GetBill() data.Bill
	PayBill(billID int, code string) bool
}

// core need to also store data and data would store requirements
