package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"saaj/core"
	"saaj/core/data"
	"saaj/input"
	"saaj/menu"
	"saaj/menu/styled"
)

type App struct {
	core         core.Core
	packageID    int
	requirements []data.Requirement
}

func NewHttpApp() *App {
	return &App{core: core.NewREST(core.Domain)}
}

func getLogin() (string, string) {
	questions := []input.Question{
		input.NewShortQuestion("please enter your username:", "username"),
		input.NewShortQuestion("please enter your password:", "password"),
	}
	main := input.NewParallelQuestion(questions)

	p := tea.NewProgram(*main, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	return questions[0].Answer, questions[1].Answer
}

func (a *App) login() {
	for {
		username, password := getLogin()
		fmt.Print("\033[2J")
		err, prompt := a.core.Authenticate(username, password)
		fmt.Println(prompt)
		if err == nil {
			break
		}
		fmt.Println("error: ", err)
	}
}

func (a *App) showPackages() {
	packages := a.core.GetPackage()
	//packages :=
	//	[]data.Package{
	//		data.Package{
	//			ID:    2,
	//			Class: "Silver",
	//		},
	//	}

	board := styled.NewBoard()
	board.InitPackageMenu([]string{"Golden", "Silver"}, packages)
	p := tea.NewProgram(board)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Print("\033[2J")
	fmt.Println("Please enter package ID:")
	var packID int
	_, err := fmt.Scanf("%d", &packID)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	a.packageID = packID

	a.requirements = a.core.RequestPackage(packID)
}

func (a *App) showRequests() {

}

func (a *App) uploadDocuments() {
	if len(a.requirements) == 0 {
		fmt.Printf("You don't need to upload any document!")
		return
	}
	questions := make([]input.Question, 0)
	for r := range a.requirements {
		nq := input.NewShortQuestion(a.requirements[r].Name+": "+a.requirements[r].Description, "enter file path")
		questions = append(questions, nq)
	}
	ps_ := input.NewSerialQuestions(questions)
	p := tea.NewProgram(*ps_)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	docs := make([]data.Document, 0)
	reqs := make([]data.Requirement, 0)
	for i := range a.requirements {
		d, err := data.NewDocument(a.requirements[i].ID, questions[i].Answer)
		if err != nil {
			fmt.Printf("couldn't catch doc(%s), error: %v\n", questions[i].Answer, err)
			reqs = append(reqs, a.requirements[i])
		}
		docs = append(docs, d)
	}
	a.requirements = reqs
	_ = a.core.SubmitDocuments(a.packageID, docs) //TODO: handle this error
}

func (a *App) reserveHotel() {

}

func (a *App) step() {
	initialModel := menu.Model{}
	initialModel.InitBaseMenu()

	p := tea.NewProgram(initialModel)
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error running program: %v\n", err)
	}

	itemName := initialModel.Items[m.(menu.Model).Selected].Name
	switch itemName {
	case "Packages":
		a.showPackages()
	case "Requests":
		a.showRequests()
	case "Upload Documents":
		a.uploadDocuments()
	case "Reserve hotel":
		a.reserveHotel()
	case "Request Visa":
		panic("Implement me!")
	case "Submit Visa":
		panic("Implement me!")
	case "Pay bill":
		panic("Implement me!")
	case "Logout":
		panic("Implement me!")
	}
}

func (a *App) Run() {
	a.step()
}
