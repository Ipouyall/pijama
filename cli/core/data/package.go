package data

type Package struct {
	ID          int    `json:"id"`
	Name        string `json:"package_name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Cost        string `json:"estimated_cost"`
	City        string `json:"city"`
	Doctor      string `json:"doctor"`
	Hospital    string `json:"hospital"`
	Class       string `json:"package_class"`
}
