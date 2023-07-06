package input

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SerialQuestion struct {
	styles    *Styles
	index     int
	questions []Question
	width     int
	height    int
	done      bool
}

func NewSerialQuestions(questions []Question) *SerialQuestion {
	styles := DefaultStyles()
	return &SerialQuestion{styles: styles, questions: questions}
}

func (m SerialQuestion) Init() tea.Cmd {
	return m.questions[m.index].input.Blink
}

func (m SerialQuestion) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	current := &m.questions[m.index]
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q": // to quit the app
			return m, tea.Quit
		case "enter": // to set the command
			if m.index == len(m.questions)-1 {
				m.done = true
			}
			current.Answer = current.input.Value()
			m.Next()
			return m, current.input.Blur
		}
	}
	current.input, cmd = current.input.Update(msg)
	return m, cmd
}

func (m SerialQuestion) View() string {
	current := m.questions[m.index]
	if m.done {
		var output string
		//for _, q := range m.questions {
		//	output += fmt.Sprintf("%s: %s\n", q.question, q.Answer)
		//}
		return output
	}
	if m.width == 0 {
		return "loading..."
	}
	// stack some left-aligned strings together in the center of the window
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Left,
			current.question,
			m.styles.InputField.Render(current.input.View()),
		),
	)
}

func (m *SerialQuestion) Next() {
	if m.index < len(m.questions)-1 {
		m.index++
		return
	}
	m.index = 0
}
