package styled

import (
	"github.com/charmbracelet/bubbles/list"
	"saaj/core/data"
)

func (b *Board) InitBaseMenu() {
	b.last_index = third
	b.cols = []column{
		newColumn(first, b.last_index),
		newColumn(second, b.last_index),
		newColumn(third, b.last_index),
	}

	b.cols[first].list.Title = "Package"
	b.cols[first].list.SetItems([]list.Item{
		NewItem(first, b.last_index, "Packages", "to see and choose a data"),
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

func (b *Board) InitPackageMenu(cols []string, packs []data.Package) {
	b.last_index = index(len(cols) - 1)
	b.cols = []column{}
	columns := make(map[string][]list.Item)
	for i := range cols {
		b.cols = append(b.cols, newColumn(index(i), b.last_index))
		b.cols[i].list.Title = cols[i]
		columns[cols[i]] = make([]list.Item, 0)
	}
	for i := range packs {
		columns[packs[i].Class] = append(columns[packs[i].Class], packs[i])
	}
	for i := range cols {
		b.cols[i].list.SetItems(columns[cols[i]])
	}
}
