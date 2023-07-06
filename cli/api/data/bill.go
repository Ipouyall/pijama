package data

type Bill struct {
	Title     string `json:"title"`
	Cost      int    `json:"total_cost"`
	PaymentID int    `json:"payment_request_id"`
}
