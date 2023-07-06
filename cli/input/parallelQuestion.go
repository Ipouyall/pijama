package input

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ParallelQuestion struct {
	styles    *Styles
	deadStyle *DeadStyle
	questions []Question
	width     int
	height    int
	done      bool
	curser    int
}

func NewParallelQuestion(questions []Question) *ParallelQuestion {
	styles := DefaultStyles()
	deadS := NewDeadStyle()
	return &ParallelQuestion{
		styles:    styles,
		deadStyle: deadS,
		questions: questions,
		curser:    0,
	}
}

func (pq ParallelQuestion) Init() tea.Cmd {
	var cmds []tea.Cmd
	for i := range pq.questions {
		cmds = append(cmds, pq.questions[i].input.Blink)
	}
	return tea.Batch(cmds...)
}

func (pq ParallelQuestion) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		pq.width = msg.Width
		pq.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc": // to quit the app
			return pq, tea.Quit
		case "enter": // to set the command
			current := &pq.questions[pq.curser]
			current.Answer = current.input.Value()
			if pq.curser == len(pq.questions)-1 {
				if current.Answer == "" {
					return pq, cmd
				}
				pq.done = true
			} else {
				if current.Answer == "" {
					return pq, cmd
				}
				pq.curser = (pq.curser + 1) % len(pq.questions)
			}
		case "up":
			current := &pq.questions[pq.curser]
			if current.input.Value() != "" {
				current.Answer = current.input.Value()
			}
			pq.curser = (pq.curser - 1 + len(pq.questions)) % len(pq.questions)
			current = &pq.questions[pq.curser]
			current.input, cmd = current.input.Update(msg)
		case "down":
			current := &pq.questions[pq.curser]
			if current.input.Value() != "" {
				current.Answer = current.input.Value()
			}
			pq.curser = (pq.curser + 1) % len(pq.questions)
			current = &pq.questions[pq.curser]
			current.input, cmd = current.input.Update(msg)
		}
	}

	if !pq.done {
		pq.questions[pq.curser].input.Update(msg)
	} else {
		return pq, tea.Quit
	}

	return pq, cmd
}

func (pq ParallelQuestion) View() string {
	if pq.done {
		var output string
		//for _, q := range pq.questions {
		//	output += fmt.Sprintf("%s: %s\n", q.question, q.Answer)
		//}
		return output
	}

	if pq.width == 0 {
		return "loading..."
	}

	var views []string
	for idx, q := range pq.questions {
		if idx == pq.curser {
			views = append(views,
				lipgloss.JoinVertical(
					lipgloss.Left,
					q.question,
					pq.styles.InputField.Render(q.input.View()),
				),
			)
			continue
		}
		views = append(views,
			lipgloss.JoinVertical(
				lipgloss.Left,
				q.question,
				pq.deadStyle.InputField.Render(q.input.View()),
			),
		)
	}

	return lipgloss.Place(
		pq.width,
		pq.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Left,
			views...,
		),
	)
}
