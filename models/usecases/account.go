package usecases

import (
	"Dp218Go/models"
	"time"
)

type AccountUsecases interface {
	CalculateMoneyAmountByDate(account models.Account, byTime time.Time) (int, error)
	CalculateProfitForPeriod(account models.Account, start, end time.Time) (int, error)
	CalculateLossForPeriod(account models.Account, start, end time.Time) (int, error)
	AddMoneyToAccount(account models.Account, amountCents int) error
	TakeMoneyFromAccount(account models.Account, amountCents int) error
}