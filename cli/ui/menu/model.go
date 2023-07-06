package menu

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"saaj/api/data"
	"strconv"
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
		case "ctrl+c", "esc", "q":
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
			selected := color.GreenString("►")
			view += fmt.Sprintf("%s %s %s\n", selected, item.Name, desc)
		} else {
			view += fmt.Sprintf("  %s %s\n", item.Name, desc)
		}
	}
	return view
}

type PackagesModel struct {
	Packages []data.Package
	Selected int
}

func NewPackageModel(packages []data.Package) *PackagesModel {
	return &PackagesModel{
		Packages: packages,
		Selected: 0,
	}
}

func (m PackagesModel) Init() tea.Cmd {
	return nil
}

func (m PackagesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
		case "up":
			m.Selected--
			if m.Selected < 0 {
				m.Selected = len(m.Packages) - 1
			}
		case "down":
			m.Selected++
			if m.Selected >= len(m.Packages) {
				m.Selected = 0
			}
		case "enter":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		// Not needed
	}

	return m, nil
}

func (m PackagesModel) View() string {
	view := ""
	for i, pkg := range m.Packages {
		if i == m.Selected {
			nameStyle := lipgloss.NewStyle().Bold(true)
			view += nameStyle.Render(color.GreenString("[x]")+"\t"+pkg.Name) + "\n"
		} else {
			view += "[ ] " + pkg.Name + "\n"
		}
		// Display other fields of the package
		view += "\tCategory: " + pkg.Category + "\n"
		view += "\tClass: " + pkg.Class + "\n"
		view += "\tCost: " + strconv.Itoa(pkg.Cost) + "\n"
		view += "\tCity: " + pkg.City + "\n"
		view += "\tDoctor: " + pkg.Doctor + "\n"
		view += "\tHospital: " + pkg.Hospital + "\n"
		view += "\tDescription: " + pkg.PDescription + "\n\n"
	}

	return view
}

type HotelRoomsModel struct {
	Rooms    []data.HotelRoom
	Selected int
}

func NewHotelRoomsModel(rooms []data.HotelRoom) *HotelRoomsModel {
	return &HotelRoomsModel{
		Rooms:    rooms,
		Selected: 0,
	}
}

func (m HotelRoomsModel) Init() tea.Cmd {
	return nil
}

func (m HotelRoomsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
		case "up":
			m.Selected--
			if m.Selected < 0 {
				m.Selected = len(m.Rooms) - 1
			}
		case "down":
			m.Selected++
			if m.Selected >= len(m.Rooms) {
				m.Selected = 0
			}
		case "enter":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		// Not needed
	}

	return m, nil
}

func (m HotelRoomsModel) View() string {
	view := ""
	for i, room := range m.Rooms {
		if i == m.Selected {
			nameStyle := lipgloss.NewStyle().Bold(true)
			view += nameStyle.Render(color.GreenString("[x]")+"\t"+room.HotelName) + "\n"
		} else {
			view += "[ ] " + room.HotelName + "\n"
		}
		view += "\tClass: " + room.HotelClass + "\n"
		view += "\tCost: " + strconv.Itoa(room.Cost) + "\n"
		view += "\tCity: " + room.City + "\n"
		view += "\tAddress: " + room.Address + "\n\n"
	}

	return view
}

type VisaStatusModel struct {
	VisaStatus []data.VisaStatus
	Selected   int
}

func NewVisaStatusModel(status []data.VisaStatus) *VisaStatusModel {
	return &VisaStatusModel{
		VisaStatus: status,
		Selected:   0,
	}
}

func (m VisaStatusModel) Init() tea.Cmd {
	return nil
}

func (m VisaStatusModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
		case "up":
			m.Selected--
			if m.Selected < 0 {
				m.Selected = len(m.VisaStatus) - 1
			}
		case "down":
			m.Selected++
			if m.Selected >= len(m.VisaStatus) {
				m.Selected = 0
			}
		case "enter":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		// Not needed
	}

	return m, nil
}

func (m VisaStatusModel) View() string {
	view := ""
	for i, status := range m.VisaStatus {
		sec := ""
		sec += "--- " + status.Serial + "\n"
		sec += "\tState: " + status.State + "\n\n"

		if i == m.Selected {
			view += lipgloss.NewStyle().Bold(true).Render(sec)
		} else {
			view += sec
		}
	}
	return view
}

type BillMenuModel struct {
	Bills    []data.Bill
	Selected int
}

func NewBillMenuModel(bills []data.Bill) *BillMenuModel {
	return &BillMenuModel{
		Bills:    bills,
		Selected: 0,
	}
}

func (m BillMenuModel) Init() tea.Cmd {
	return nil
}

func (m BillMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "up":
			m.Selected--
			if m.Selected < 0 {
				m.Selected = len(m.Bills) - 1
			}
		case "down":
			m.Selected++
			if m.Selected >= len(m.Bills) {
				m.Selected = 0
			}
		case "enter":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m BillMenuModel) View() string {
	view := ""
	for i, bill := range m.Bills {
		if i == m.Selected {
			// Bold style for the selected bill
			view += lipgloss.NewStyle().Bold(true).Render("[x] "+bill.Title) + "\n"
		} else {
			view += "[ ] " + bill.Title + "\n"
		}
		view += "\tCost: $" + strconv.Itoa(bill.Cost) + "\n"
		view += "\tPayment ID: " + strconv.Itoa(bill.PaymentID) + "\n\n"
	}

	return view
}
