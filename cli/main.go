package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
	"saaj/input"
	"saaj/menu"
)

func login() {
	questions := []input.Question{
		input.NewShortQuestion("please enter your username:", "username"),
		input.NewShortQuestion("please enter your password:", "password"),
	}
	main := input.NewParallelQuestion(questions)

	p := tea.NewProgram(*main, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	login()

	board := menu.NewBoard()
	board.InitLists()
	p := tea.NewProgram(board)
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
