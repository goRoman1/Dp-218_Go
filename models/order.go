package models

type Order struct {
	ID            int     `json:"id"`
	UserID        int     `json:"user_id"`
	ScooterID     int     `json:"scooter_id"`
	StatusStartID int     `json:"status_start_id"`
	StatusEndID   int     `json:"status_end_id"`
	Distance      float64 `json:"distance"`
	Amount        int 	  `json:"amount"`
}

type OrderList struct {
	Orders []Order `json:"orders"`
}
