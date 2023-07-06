package core

import (
	"saaj/core/data"
)

type Core interface {
	Authenticate(username, password string) (error, string)

	GetPackage() []data.Package
	RequestPackage(pack data.Package) []data.Requirement

	SubmitDocuments(packID int, docs []data.Document) error

	GetHotels() []data.HotelRoom
	ReserveHotel(hotelID int) error

	RequestVisa() []data.Requirement
	SubmitVisa(visaID int) error

	GetBill() data.Bill
	PayBill(billID int, code string) error
}

// TODO: needs to see its requests + logout
