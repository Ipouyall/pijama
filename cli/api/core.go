package api

import (
	"saaj/api/data"
)

type Core interface {
	Authenticate(username, password string) (error, string)

	GetPackage() []data.Package
	RequestPackage(pack data.Package) []data.Requirement

	SubmitDocuments(packID int, docs []data.Document, dKind string) (error, data.Bill)

	GetHotels() []data.HotelRoom
	ReserveHotel(hotelID int) error

	RequestVisa() []data.Requirement
	VisaStatus() []data.VisaStatus

	GetBill() data.Bill
	PayBill(billID int, code string) error

	Logout() error
}

// TODO: needs to see its requests + logout
