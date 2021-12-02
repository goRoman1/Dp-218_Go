package postgres

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
	"context"
	"strconv"
	"strings"
	"time"
)

type AccountRepoDB struct {
	userRepo *UserRepoDB
	db       repositories.AnyDatabase
}

func NewAccountRepoDB(userRepo *UserRepoDB, db repositories.AnyDatabase) *AccountRepoDB {
	return &AccountRepoDB{userRepo, db}
}

func (accdb *AccountRepoDB) GetAccountsByOwner(user models.User) (*models.AccountList, error) {
	list := &models.AccountList{}

	querySQL := `SELECT id, name, number FROM accounts WHERE owner_id = $1;`
	rows, err := accdb.db.QueryResult(context.Background(), querySQL, user.ID)
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var account models.Account
		err := rows.Scan(&account.ID, &account.Name, &account.Number)
		if err != nil {
			return list, err
		}
		account.User = user
		list.Accounts = append(list.Accounts, account)
	}

	return list, nil
}

func (accdb *AccountRepoDB) GetAccountById(accountId int) (models.Account, error) {
	account := models.Account{}

	querySQL := `SELECT id, name, number, owner_id FROM accounts WHERE id = $1;`
	row := accdb.db.QueryResultRow(context.Background(), querySQL, accountId)
	var userId int
	err := row.Scan(&account.ID, &account.Name, &account.Number, &userId)
	account.User, err = accdb.userRepo.GetUserById(userId)

	return account, err
}

func (accdb *AccountRepoDB) GetAccountByNumber(number string) (models.Account, error) {
	account := models.Account{}

	querySQL := `SELECT id, name, number, owner_id FROM accounts WHERE number = $1;`
	row := accdb.db.QueryResultRow(context.Background(), querySQL, number)
	var userId int
	err := row.Scan(&account.ID, &account.Name, &account.Number, &userId)
	account.User, err = accdb.userRepo.GetUserById(userId)

	return account, err
}

func (accdb *AccountRepoDB) AddAccount(account *models.Account) error {
	var id int
	querySQL := `INSERT INTO accounts(name, number, owner_id) VALUES($1, $2, $3) RETURNING id;`
	err := accdb.db.QueryResultRow(context.Background(), querySQL, account.Name, account.Number, account.User.ID).Scan(&id)
	if err != nil {
		return err
	}
	account.ID = id

	return nil
}

func (accdb *AccountRepoDB) UpdateAccount(accountId int, accountData models.Account) (models.Account, error) {
	account := models.Account{}
	querySQL := `UPDATE accounts SET name=$1, number=$2, owner_id=$3 WHERE id=$4 RETURNING id, name, number, owner_id;`
	var userId int
	err := accdb.db.QueryResultRow(context.Background(), querySQL, account.Name, account.Number, account.User.ID).Scan(&account.ID, &account.Name, &account.Number, &userId)
	if err != nil {
		return account, err
	}
	account.User, err = accdb.userRepo.GetUserById(userId)
	if err != nil {
		return account, err
	}

	return account, nil
}

func (accdb *AccountRepoDB) GetAccountTransactionById(transId int) (models.AccountTransaction, error) {
	accountTransaction := models.AccountTransaction{}

	querySQL := `SELECT id, date_time, payment_type_id, account_from_id, account_to_id, order_id, amount_cents FROM account_transactions WHERE id = $1;`
	row := accdb.db.QueryResultRow(context.Background(), querySQL, transId)
	var paymentId int
	var accFromId, accToId int
	var orderId int
	err := row.Scan(&accountTransaction.ID, &accountTransaction.DateTime, &paymentId, &accFromId, &accToId, orderId, &accountTransaction.AmountCents)
	if err != nil {
		return accountTransaction, err
	}

	err = addTransactionComplexFields(accdb, &accountTransaction, paymentId, accFromId, accToId, orderId)
	if err != nil {
		return accountTransaction, err
	}

	return accountTransaction, err
}

func addTransactionComplexFields(accdb *AccountRepoDB, accountTransaction *models.AccountTransaction, paymentId, accFromId, accToId, orderId int) error {
	var err error
	accountTransaction.PaymentType, err = accdb.GetPaymentTypeById(paymentId)
	if err != nil {
		return err
	}
	accountTransaction.AccountFrom, err = accdb.GetAccountById(accFromId)
	if err != nil && accFromId != 0 {
		return err
	}
	accountTransaction.AccountTo, err = accdb.GetAccountById(accToId)
	if err != nil && accToId != 0 {
		return err
	}
	accountTransaction.Order, err = models.Order{}, nil //TODO: refactor when Orders implemented
	if err != nil && orderId != 0 {
		return err
	}
	return nil
}

func (accdb *AccountRepoDB) AddAccountTransaction(accountTransaction *models.AccountTransaction) error {
	var id int
	querySQL := `INSERT INTO account_transactions(date_time, payment_type_id, account_from_id, account_to_id, order_id, amount_cents) VALUES($1, $2, $3, $4, $5, $6) RETURNING id;`
	err := accdb.db.QueryResultRow(context.Background(), querySQL, accountTransaction.DateTime, accountTransaction.PaymentType.ID,
		accountTransaction.AccountFrom.ID, accountTransaction.AccountTo.ID, accountTransaction.Order.ID, accountTransaction.AmountCents).Scan(&id)
	if err != nil {
		return err
	}
	accountTransaction.ID = id

	return nil
}

func getTransactionsBySomeQuery(accdb *AccountRepoDB, querySQL string, params ...interface{}) (*models.AccountTransactionList, error) {
	list := &models.AccountTransactionList{}
	rows, err := accdb.db.QueryResult(context.Background(), querySQL, params)
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var accountTransaction models.AccountTransaction
		var paymentId int
		var accFromId, accToId int
		var orderId int
		err := rows.Scan(&accountTransaction.ID, &accountTransaction.DateTime, &paymentId, &accFromId, &accToId, &orderId, &accountTransaction.AmountCents)
		if err != nil {
			return list, err
		}

		err = addTransactionComplexFields(accdb, &accountTransaction, paymentId, accFromId, accToId, orderId)
		if err != nil {
			return list, err
		}

		list.AccountTransactions = append(list.AccountTransactions, accountTransaction)
	}

	return list, nil
}

func (accdb *AccountRepoDB) GetAccountTransactions(accounts ...models.Account) (*models.AccountTransactionList, error) {
	querySQL := `SELECT id, date_time, payment_type_id, account_from_id, account_to_id, order_id, amount_cents FROM account_transactions`
	var params []int
	for i, acc := range accounts {
		if i == 0 {
			querySQL += ` WHERE FALSE`
		}
		paramIndex := strconv.Itoa(i + 1)
		querySQL += ` OR account_from_id = $` + paramIndex + ` OR account_to_id = $` + paramIndex
		params = append(params, acc.ID)
	}
	querySQL += `;`

	return getTransactionsBySomeQuery(accdb, querySQL, params)
}

func (accdb *AccountRepoDB) GetAccountTransactionsInTimePeriod(start time.Time, end time.Time, accounts ...models.Account) (*models.AccountTransactionList, error) {
	querySQL := `SELECT id, date_time, payment_type_id, account_from_id, account_to_id, order_id, amount_cents FROM account_transactions
		WHERE date_time>=$1 AND date_time<=$2`
	var params []interface{}
	params = append(params, start)
	params = append(params, end)
	var accountCondition = make([]string, len(accounts))
	for i, acc := range accounts {
		accountCondition[i] = `$` + strconv.Itoa(i+3)
		params = append(params, acc.ID)
	}
	if len(accounts) > 0 {
		conditionStr := strings.Join(accountCondition, ", ")
		querySQL += ` AND (account_from_id IN (` + conditionStr + `) OR account_to_id IN (` + conditionStr + `))`
	}
	querySQL += `;`

	return getTransactionsBySomeQuery(accdb, querySQL, params)
}

func (accdb *AccountRepoDB) GetAccountTransactionsByOrder(order models.Order) (*models.AccountTransactionList, error) {
	querySQL := `SELECT id, date_time, payment_type_id, account_from_id, account_to_id, order_id, amount_cents FROM account_transactions
		WHERE order_id=$1;`

	return getTransactionsBySomeQuery(accdb, querySQL, order.ID)
}

func (accdb *AccountRepoDB) GetAccountTransactionsByPaymentType(paymentType models.PaymentType, accounts ...models.Account) (*models.AccountTransactionList, error){
	querySQL := `SELECT id, date_time, payment_type_id, account_from_id, account_to_id, order_id, amount_cents FROM account_transactions
		WHERE payment_type_id=$1`
	var params []int
	params = append(params, paymentType.ID)
	var accountCondition = make([]string, len(accounts))
	for i, acc := range accounts {
		accountCondition[i] = `$` + strconv.Itoa(i+2)
		params = append(params, acc.ID)
	}
	if len(accounts) > 0 {
		conditionStr := strings.Join(accountCondition, ", ")
		querySQL += ` AND (account_from_id IN (` + conditionStr + `) OR account_to_id IN (` + conditionStr + `))`
	}
	querySQL += `;`

	return getTransactionsBySomeQuery(accdb, querySQL, params)
}

func (accdb *AccountRepoDB) GetPaymentTypeById(paymentTypeId int) (models.PaymentType, error) {
	paymentType := models.PaymentType{}

	querySQL := `SELECT id, name FROM payment_types WHERE id = $1;`
	row := accdb.db.QueryResultRow(context.Background(), querySQL, paymentTypeId)
	err := row.Scan(&paymentType.ID, &paymentType.Name)

	return paymentType, err
}
