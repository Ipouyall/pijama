package menu

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Board struct {
	help       help.Model
	loaded     bool
	focused    index
	cols       []column
	quitting   bool
	last_index index
}

func NewBoard() *Board {
	help := help.New()
	help.ShowAll = true
	return &Board{help: help, focused: first}
}

func (m *Board) Init() tea.Cmd {
	return nil
}

func (m *Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		var cmds []tea.Cmd
		m.help.Width = msg.Width - margin
		for i := 0; i < len(m.cols); i++ {
			var res tea.Model
			res, cmd = m.cols[i].Update(msg)
			m.cols[i] = res.(column)
			cmds = append(cmds, cmd)
		}
		m.loaded = true
		return m, tea.Batch(cmds...)
	case Form:
		return m, m.cols[m.focused].Set(msg.index, msg.CreateTask())
	case moveMsg:
		return m, m.cols[m.focused.getNext(m.last_index)].Set(APPEND, msg.Item)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, keys.Left):
			m.cols[m.focused].Blur()
			m.focused = m.focused.getPrev(m.last_index)
			m.cols[m.focused].Focus()
		case key.Matches(msg, keys.Right):
			m.cols[m.focused].Blur()
			m.focused = m.focused.getNext(m.last_index)
			m.cols[m.focused].Focus()
		}
	}
	res, cmd := m.cols[m.focused].Update(msg)
	if _, ok := res.(column); ok {
		m.cols[m.focused] = res.(column)
	} else {
		return res, cmd
	}
	return m, cmd
}

// Changing to pointer receiver to get back to this model after adding a new task via the form... Otherwise I would need to pass this model along to the form and it becomes highly coupled to the other models.
func (m *Board) View() string {
	if m.quitting {
		return ""
	}
	if !m.loaded {
		return "loading..."
	}
	var views []string
	for idx := range m.cols {
		views = append(views, m.cols[idx].View())
	}
	board := lipgloss.JoinHorizontal(
		lipgloss.Left,
		views...,
	)
	return lipgloss.JoinVertical(lipgloss.Left, board, m.help.View(keys))
}
