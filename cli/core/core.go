package core

import _package "saaj/package"

type Core interface {
	Authenticate(username, password string) (bool, string)

	GetPackage() []_package.Package
	RequestPackage(packID int) _package.Requirements

	SubmitDocument(docID int, name, content string) bool

	GetHotels() []_package.HotelRoom
	ReserveHotel(hotelID int) bool

	RequestVisa() bool
	SubmitVisa(visaID int) bool

	GetBill() _package.Bill
	PayBill(billID int, code string) bool
}

// core need to also store package and package would store requirements
