package menu

import "github.com/charmbracelet/bubbles/list"

func (b *Board) InitBaseMenu() {
	b.last_index = third
	b.cols = []column{
		newColumn(first, b.last_index),
		newColumn(second, b.last_index),
		newColumn(third, b.last_index),
	}

	b.cols[first].list.Title = "Package"
	b.cols[first].list.SetItems([]list.Item{
		NewItem(first, b.last_index, "Packages", "to see and choose a package"),
		NewItem(first, b.last_index, "Upload Documents", "upload your docs"),
		NewItem(first, b.last_index, "Reserve hotel", "reserve a hotel"),
	})

	b.cols[second].list.Title = "Visa"
	b.cols[second].list.SetItems([]list.Item{
		NewItem(first, b.last_index, "Request Visa", "if you need visa"),
		NewItem(first, b.last_index, "Submit Visa", "if you have visa"),
	})

	b.cols[third].list.Title = "Dashboard"
	b.cols[third].list.SetItems([]list.Item{
		NewItem(first, b.last_index, "Requests", "check you requests status"),
		NewItem(first, b.last_index, "Pay bill", "make a payment"),
		NewItem(first, b.last_index, "Logout", "to logout"),
	})
}
