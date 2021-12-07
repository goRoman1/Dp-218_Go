package services

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
	"time"
)

const (
	PayIncomeTypeId  = 2
	PayOutcomeTypeId = 3
)

type AccountService struct {
	repoAccount            repositories.AccountRepo
	repoAccountTransaction repositories.AccountTransactionRepo
	repoPaymentType        repositories.PaymentTypeRepo
}

type transactionsWithIncome struct {
	Transaction models.AccountTransaction
	IsIncome    bool
}

func NewAccountService(repoAccount repositories.AccountRepo, repoAccountTransaction repositories.AccountTransactionRepo, repoPaymentType repositories.PaymentTypeRepo) *AccountService {
	return &AccountService{repoAccount, repoAccountTransaction, repoPaymentType}
}

func (accserv *AccountService) GetAccountsByOwner(user models.User) (*models.AccountList, error) {
	return accserv.repoAccount.GetAccountsByOwner(user)
}

func (accserv *AccountService) GetAccountById(accountId int) (models.Account, error) {
	return accserv.repoAccount.GetAccountById(accountId)
}

func (accserv *AccountService) GetAccountByNumber(number string) (models.Account, error) {
	return accserv.repoAccount.GetAccountByNumber(number)
}

func (accserv *AccountService) AddAccount(account *models.Account) error {
	return accserv.repoAccount.AddAccount(account)
}

func (accserv *AccountService) UpdateAccount(accountId int, accountData models.Account) (models.Account, error) {
	return accserv.repoAccount.UpdateAccount(accountId, accountData)
}

func (accserv *AccountService) GetAccountTransactionById(transId int) (models.AccountTransaction, error) {
	return accserv.repoAccountTransaction.GetAccountTransactionById(transId)
}

func (accserv *AccountService) AddAccountTransaction(accountTransaction *models.AccountTransaction) error {
	return accserv.repoAccountTransaction.AddAccountTransaction(accountTransaction)
}

func (accserv *AccountService) GetAccountTransactions(accounts ...models.Account) (*models.AccountTransactionList, error) {
	return accserv.repoAccountTransaction.GetAccountTransactions(accounts...)
}

func (accserv *AccountService) GetAccountTransactionsInTimePeriod(start time.Time, end time.Time, accounts ...models.Account) (*models.AccountTransactionList, error) {
	return accserv.repoAccountTransaction.GetAccountTransactionsInTimePeriod(start, end, accounts...)
}

func (accserv *AccountService) GetAccountTransactionsByOrder(order models.Order) (*models.AccountTransactionList, error) {
	return accserv.repoAccountTransaction.GetAccountTransactionsByOrder(order)
}

func (accserv *AccountService) GetAccountTransactionsByPaymentType(paymentType models.PaymentType, accounts ...models.Account) (*models.AccountTransactionList, error) {
	return accserv.repoAccountTransaction.GetAccountTransactionsByPaymentType(paymentType, accounts...)
}

func (accserv *AccountService) GetPaymentTypeById(paymentTypeId int) (models.PaymentType, error) {
	return accserv.repoPaymentType.GetPaymentTypeById(paymentTypeId)
}

func (accserv *AccountService) CalculateMoneyAmountByDate(account models.Account, byTime time.Time) (models.Money, error) {
	transactionsUpToDate, err := accserv.repoAccountTransaction.GetAccountTransactionsInTimePeriod(time.UnixMilli(0), byTime, account)
	if err != nil {
		return models.Money{}, err
	}
	var amountCalculated int
	for _, trans := range transactionsUpToDate.AccountTransactions {
		if trans.AccountFrom.ID == account.ID {
			amountCalculated -= trans.AmountCents
		}
		if trans.AccountTo.ID == account.ID {
			amountCalculated += trans.AmountCents
		}
	}
	return accserv.MoneyFromCents(amountCalculated), nil
}

func (accserv *AccountService) CalculateProfitForPeriod(account models.Account, start, end time.Time) (models.Money, error) {
	transactionsUpToDate, err := accserv.repoAccountTransaction.GetAccountTransactionsInTimePeriod(start, end, account)
	if err != nil {
		return models.Money{}, err
	}
	var amountCalculated int
	for _, trans := range transactionsUpToDate.AccountTransactions {
		if trans.AccountTo.ID == account.ID {
			amountCalculated += trans.AmountCents
		}
	}
	return accserv.MoneyFromCents(amountCalculated), nil
}

func (accserv *AccountService) CalculateLossForPeriod(account models.Account, start, end time.Time) (models.Money, error) {
	transactionsUpToDate, err := accserv.repoAccountTransaction.GetAccountTransactionsInTimePeriod(start, end, account)
	if err != nil {
		return models.Money{}, err
	}
	var amountCalculated int
	for _, trans := range transactionsUpToDate.AccountTransactions {
		if trans.AccountFrom.ID == account.ID {
			amountCalculated += trans.AmountCents
		}
	}
	return accserv.MoneyFromCents(amountCalculated), nil
}

func (accserv *AccountService) AddMoneyToAccount(account models.Account, amountCents int) error {
	paymentType, err := accserv.repoPaymentType.GetPaymentTypeById(PayIncomeTypeId)
	if err != nil {
		return err
	}

	accTransaction := &models.AccountTransaction{
		DateTime:    time.Now(),
		PaymentType: paymentType,
		AccountFrom: models.Account{},
		AccountTo:   account,
		Order:       models.Order{},
		AmountCents: amountCents}

	return accserv.repoAccountTransaction.AddAccountTransaction(accTransaction)
}

func (accserv *AccountService) TakeMoneyFromAccount(account models.Account, amountCents int) error {
	paymentType, err := accserv.repoPaymentType.GetPaymentTypeById(PayOutcomeTypeId)
	if err != nil {
		return err
	}

	accTransaction := &models.AccountTransaction{
		DateTime:    time.Now(),
		PaymentType: paymentType,
		AccountFrom: account,
		AccountTo:   models.Account{},
		Order:       models.Order{},
		AmountCents: amountCents}

	return accserv.repoAccountTransaction.AddAccountTransaction(accTransaction)
}

func (accserv *AccountService) MoneyFromCents(cents int) models.Money {
	return models.Money{
		Dollars: cents / 100,
		Cents:   cents % 100,
	}
}

func (accserv *AccountService) CentsFromMoney(money models.Money) int {
	return money.Dollars*100 + money.Cents
}

func (accserv *AccountService) GetAccountOutputStructById(accId int) (interface{}, error) {
	account, err := accserv.GetAccountById(accId)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	moneyTotal, err := accserv.CalculateMoneyAmountByDate(account, now)
	if err != nil {
		return nil, err
	}

	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthIncome, err := accserv.CalculateProfitForPeriod(account, monthStart, now)
	if err != nil {
		return nil, err
	}
	monthOutcome, err := accserv.CalculateLossForPeriod(account, monthStart, now)
	if err != nil {
		return nil, err
	}
	monthTransactions, err := accserv.GetAccountTransactionsInTimePeriod(monthStart, now, account)
	if err != nil {
		return nil, err
	}

	return struct {
		ID                  int
		Number              string
		Name                string
		TotalAmount         models.Money
		MonthlyIncome       models.Money
		MonthlyOutcome      models.Money
		MonthlyTransactions []transactionsWithIncome
	}{
		ID:                  account.ID,
		Number:              account.Number,
		Name:                account.Name,
		TotalAmount:         moneyTotal,
		MonthlyIncome:       monthIncome,
		MonthlyOutcome:      monthOutcome,
		MonthlyTransactions: addIncomeToTransactions(monthTransactions.AccountTransactions, account),
	}, nil
}

func addIncomeToTransactions(transactions []models.AccountTransaction, account models.Account) []transactionsWithIncome {
	result := make([]transactionsWithIncome, len(transactions))
	for i := 0; i < len(transactions); i++ {
		result[i].Transaction = transactions[i]
		result[i].IsIncome = account.ID == transactions[i].AccountTo.ID
	}
	return result
}
