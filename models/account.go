package models

import "time"

type PaymentType struct{
	ID     int    `json:"id"`
	Name   string `json:"name"`
}

type Account struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Number string `json:"is_admin"`
	User   User   `json:"user"`
}

type AccountList struct {
	Accounts []Account `json:"accounts"`
}


type AccountTransaction struct {
	ID          int         `json:"id"`
	DateTime    time.Time   `json:"date_time"`
	PaymentType PaymentType `json:"payment_type"`
	AccountFrom Account     `json:"account_from"`
	AccountTo   Account     `json:"account_to"`
	Order       Order       `json:"order"`
	AmountCents int        `json:"amount_cents"`
}

type AccountTransactionList struct {
	AccountTransactions []AccountTransaction `json:"account_transactions"`
}
