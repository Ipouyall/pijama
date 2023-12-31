package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
	"os/exec"
	"saaj/api"
	"saaj/api/data"
	input2 "saaj/ui/input"
	"saaj/ui/menu"
	"strconv"
	"time"
)

type App struct {
	core         api.Core
	packageID    int
	requirements []data.Requirement
	visaProc     bool
	Bills        map[string]data.Bill
}

func NewHttpApp() *App {
	return &App{
		core:     api.NewREST(api.Domain),
		visaProc: false,
		Bills:    make(map[string]data.Bill),
	}
}

func getLogin() (string, string) {
	questions := []input2.Question{
		input2.NewShortQuestion("please enter your username:", "username"),
		input2.NewShortQuestion("please enter your password:", "password"),
	}
	main := input2.NewParallelQuestion(questions)

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
	if a.visaProc {
		fmt.Println("You need to upload some documents for your previous request (visa) first!")
		return
	}
	packages := a.core.GetPackage()
	//packages :=
	//	[]data.Package{
	//		data.Package{
	//			ID:       2,
	//			Class:    "Silver",
	//			Name:     "bb",
	//			Category: "aas",
	//		},
	//		data.Package{
	//			ID:       5,
	//			Class:    "gerals",
	//			Name:     "bb",
	//			Category: "aas",
	//		},
	//	}

	board := menu.NewPackageModel(packages)
	p := tea.NewProgram(board)
	m, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
	pack := packages[m.(menu.PackagesModel).Selected]
	a.packageID = pack.ID

	a.requirements = a.core.RequestPackage(pack)
	fmt.Println("please upload your documents for treatment request in Upload section")
}

func (a *App) showRequests() {

}

func (a *App) uploadDocuments() {
	if len(a.requirements) == 0 {
		fmt.Printf("You don't need to upload any document!")
		return
	}
	questions := make([]input2.Question, 0)
	for r := range a.requirements {
		nq := input2.NewShortQuestion(a.requirements[r].Name+": "+a.requirements[r].Description, "enter file path")
		questions = append(questions, nq)
	}
	ps_ := input2.NewSerialQuestions(questions)
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
	var kind string
	if a.visaProc {
		kind = "Visa"
	} else {
	    kind = "Treat"
	}
	_, bill := a.core.SubmitDocuments(a.packageID, docs, kind) //TODO: handle this error
	if a.visaProc {
		a.visaProc = false
		a.Bills["visa_"+strconv.Itoa(bill.PaymentID)] = bill
	}
}

func (a *App) reserveHotel() {
	rooms := a.core.GetHotels()

	board := menu.NewHotelRoomsModel(rooms)
	p := tea.NewProgram(board)
	m, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
	room := rooms[m.(menu.HotelRoomsModel).Selected]
	_ = a.core.ReserveHotel(room.ID) //TODO: handle this error
	fmt.Println("Your room is reserved!")
}

func (a *App) requestVisa() {
	if len(a.requirements) != 0 {
		fmt.Println("You need to upload some documents for your previous requests first!")
		return
	}
	a.requirements = a.core.RequestVisa()
	a.visaProc = true
	fmt.Println("please upload your documents for visa request in Upload section")
}

func (a *App) visaStatus() {
	status := a.core.VisaStatus()

	board := menu.NewVisaStatusModel(status)
	p := tea.NewProgram(board)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (a *App) billing() {
	bill := a.core.GetBill()
	//bill := data.Bill{
	//	PaymentID: 1,
	//	Title:     "Visa",
	//	Cost:      100,
	//}

	a.Bills["treat_"+strconv.Itoa(bill.PaymentID)] = bill
	bills := make([]data.Bill, 0)
	for _, v := range a.Bills {
		bills = append(bills, v)
	}

	if len(bills) == 0 {
		fmt.Println("You don't have any bills!")
		return
	}

	board := menu.NewBillMenuModel(bills)
	p := tea.NewProgram(board)
	m, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
	billed := bills[m.(menu.BillMenuModel).Selected]
	if billed.Title == "Visa" {
		key := "visa_" + strconv.Itoa(billed.PaymentID)
		delete(a.Bills, key)
	} else if billed.Title[:5] == "Treat" {
		key := "treat_" + strconv.Itoa(billed.PaymentID)
		delete(a.Bills, key)
	}

	q := input2.NewShortQuestion(fmt.Sprintf("Enter your card number: (cost: %d)", billed.Cost), "credit-card number")
	form := input2.NewSerialQuestions([]input2.Question{q})
	p = tea.NewProgram(*form)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Print("\033[2J")
	fmt.Println("Your payment is done successfully!")
}

func (a *App) logout() {
	err := a.core.Logout()
	if err != nil {
		fmt.Println("Failed, error: ", err)
		return
	}
	fmt.Println("You logged out successfully!")
}

func (a *App) step() bool {
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
		a.requestVisa()
	case "Visa Status":
		a.visaStatus()
	case "Pay bill":
		a.billing()
	case "Logout":
		a.logout()
		return false
	}
	return true
}

func (a *App) Run() {
	a.login()

	for {
		fmt.Print("\033[2J")
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
		if a.step() == false {
			break
		}
		time.Sleep(1*time.Second)
	}
}
