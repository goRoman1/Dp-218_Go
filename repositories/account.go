package repositories

import (
	"Dp218Go/models"
	"time"
)

type AccountRepo interface {
	GetAccountsByOwner(user models.User) (*models.AccountList, error)
	GetAccountById(accountId int) (models.Account, error)
	GetAccountByNumber(number string) (models.Account, error)
	AddAccount(account *models.Account) error
	UpdateAccount(accountId int, accountData models.Account) (models.Account, error)
}

type AccountTransactionRepo interface {
	GetAccountTransactionById(transId int) (models.AccountTransaction, error)
	AddAccountTransaction(accountTransaction *models.AccountTransaction) error
	GetAccountTransactions(accounts ...models.Account) (*models.AccountTransactionList, error)
	GetAccountTransactionsInTimePeriod(start time.Time, end time.Time, accounts ...models.Account) (*models.AccountTransactionList, error)
	GetAccountTransactionsByOrder(order models.Order)(*models.AccountTransactionList, error)
	GetAccountTransactionsByPaymentType(paymentType models.PaymentType, accounts ... models.Account) (*models.AccountTransactionList, error)
}

type PaymentTypeRepo interface {
	GetPaymentTypeById(paymentTypeId int) (models.PaymentType, error)
}