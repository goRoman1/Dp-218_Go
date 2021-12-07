package models

import "time"

type PaymentType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Money struct {
	Dollars int `json:"dollars"`
	Cents   int `json:"cents"`
}

type Account struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Number string `json:"number"`
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
	AmountCents int         `json:"amount_cents"`
}

type AccountTransactionList struct {
	AccountTransactions []AccountTransaction `json:"account_transactions"`
}

func (accTrans *AccountTransaction) GetAmountInMoney() Money {
	return Money{
		Dollars: accTrans.AmountCents / 100,
		Cents:   accTrans.AmountCents % 100,
	}
}
