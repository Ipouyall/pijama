package input

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("#CF9FFF") // #CCCCFF #C3B1E1  #E0B0FF  #CF9FFF
	s.InputField = lipgloss.NewStyle().
		BorderForeground(s.BorderColor).
		BorderStyle(lipgloss.NormalBorder()).
		Padding(1).
		Width(80)
	return s
}

type DeadStyle struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

func NewDeadStyle() *DeadStyle {
	s := new(DeadStyle)
	s.BorderColor = lipgloss.Color("#130c25") // Use your desired color for not selected fields
	s.InputField = lipgloss.NewStyle().
		BorderForeground(s.BorderColor).
		BorderStyle(lipgloss.NormalBorder()).
		Padding(0).
		Width(80)
	return s
}
