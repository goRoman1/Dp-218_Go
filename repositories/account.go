package repositories

import (
	"Dp218Go/models"
	"time"
)

type AccountRepo interface {
	GetAccountsByOwner(user models.User) (*models.AccountList, error)
	GetAccountByID(accountID int) (models.Account, error)
	GetAccountByNumber(number string) (models.Account, error)
	AddAccount(account *models.Account) error
	UpdateAccount(accountID int, accountData models.Account) (models.Account, error)
}

type AccountTransactionRepo interface {
	GetAccountTransactionByID(transID int) (models.AccountTransaction, error)
	AddAccountTransaction(accountTransaction *models.AccountTransaction) error
	GetAccountTransactions(accounts ...models.Account) (*models.AccountTransactionList, error)
	GetAccountTransactionsInTimePeriod(start time.Time, end time.Time, accounts ...models.Account) (*models.AccountTransactionList, error) //nolint:lll
	GetAccountTransactionsByOrder(order models.Order) (*models.AccountTransactionList, error)
	GetAccountTransactionsByPaymentType(paymentType models.PaymentType, accounts ...models.Account) (*models.AccountTransactionList, error) //nolint:lll
}

type PaymentTypeRepo interface {
	GetPaymentTypeById(paymentTypeID int) (models.PaymentType, error)
}
