package menu

func (m *Model) InitBaseMenu() {
	items := []Item{
		{ID: 1, Name: "Packages", Description: "to show packages"},
		//{ID: 2, Name: "Requests", Description: "to show active requests"},
		{ID: 3, Name: "Upload Documents", Description: "to upload requested docs"},
		{ID: 4, Name: "Reserve hotel", Description: "to reserve a hotel"},
		{ID: 5, Name: "Request Visa", Description: "if you need visa"},
		{ID: 6, Name: "Visa Status", Description: "to see visa(s) status"},
		{ID: 7, Name: "Pay bill", Description: "make a payment"},
		{ID: 8, Name: "Logout", Description: "to logout"},
	}
	m.Items = items
}
