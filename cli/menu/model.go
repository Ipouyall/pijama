package menu

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

type Model struct {
	Items    []Item
	Selected int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "up":
			m.Selected--
			if m.Selected < 0 {
				m.Selected = len(m.Items) - 1
			}
		case "down":
			m.Selected++
			if m.Selected >= len(m.Items) {
				m.Selected = 0
			}
		case "enter":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	view := ""
	for i, item := range m.Items {
		desc := color.HiBlackString(item.Description)
		if i == m.Selected {
			selected := color.GreenString("[x]")
			view += fmt.Sprintf("%s %s %s\n", selected, item.Name, desc)
		} else {
			view += fmt.Sprintf("[ ] %s %s\n", item.Name, desc)
		}
	}
	return view
}
