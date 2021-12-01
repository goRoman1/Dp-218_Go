package usecases

import (
	"Dp218Go/models"
	"time"
)

type AccountUsecases interface {
	CalculateMoneyAmountByDate(accont models.Account, byTime time.Time) error
	AddMoneyToAccount(account models.Account, amountCents uint) error
	TakeMoneyFromAccount(account models.Account, amountCents uint) error
}
