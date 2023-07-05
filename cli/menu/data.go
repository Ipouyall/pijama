package menu

import "github.com/charmbracelet/bubbles/list"

// Provides the mock data to fill the kanban board

func (b *Board) InitLists() {
	b.cols = []column{
		newColumn(todo),
		newColumn(inProgress),
		newColumn(done),
	}
	// Init To Do
	b.cols[todo].list.Title = "Golden"
	b.cols[todo].list.SetItems([]list.Item{
		Item{status: todo, title: "buy milk", description: "strawberry milk"},
		Item{status: todo, title: "eat sushi", description: "negitoro roll, miso soup, rice"},
		Item{status: todo, title: "fold laundry", description: "or wear wrinkly t-shirts"},
	})
	// Init in progress
	b.cols[inProgress].list.Title = "Silver"
	b.cols[inProgress].list.SetItems([]list.Item{
		Item{status: inProgress, title: "write code", description: "don't worry, it's Go"},
	})
	// Init done
	b.cols[done].list.Title = "Bronze"
	b.cols[done].list.SetItems([]list.Item{
		Item{status: done, title: "stay cool", description: "as a cucumber"},
	})
}
