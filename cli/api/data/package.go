package data

import (
	"fmt"
)

type Package struct {
	ID           int    `json:"id"`
	Name         string `json:"package_name"`
	Category     string `json:"category"`
	PDescription string `json:"description"`
	Cost         int    `json:"estimated_cost"`
	City         string `json:"city"`
	Doctor       string `json:"doctor"`
	Hospital     string `json:"hospital"`
	Class        string `json:"package_class"`
}

func (t Package) FilterValue() string {
	return t.Category
}

func (t Package) Title() string {
	return t.Name
}

func (t Package) Description() string {
	desc := t.Category + "\n"
	desc += fmt.Sprintf("pid: %d\n", t.ID)
	desc += fmt.Sprintf("Cost: %d\n", t.Cost)
	desc += "Class" + t.Class + "\n"
	desc += "City: " + t.City + "\n"
	desc += "Doctor: " + t.Doctor + "\n"
	desc += "Hospital: " + t.Hospital + "\n"
	desc += "Description: " + t.PDescription + "\n"
	return desc
}
