package core

import _package "saaj/package"

type Core interface {
	Authenticate(username, password string) (bool, string)
	GetPackage() []_package.Package
	ReservePackage(pid int) bool
}
