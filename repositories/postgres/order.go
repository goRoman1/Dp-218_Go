package postgres

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
	"context"
)

type OrderRepoDb struct {
	db repositories.AnyDatabase
}

func NewOrderRepoDB(db repositories.AnyDatabase) *OrderRepoDb {
	return &OrderRepoDb{db}
}

func (ordb *OrderRepoDb) CreateOrder(user models.User, scooter models.Scooter) (models.Order, error) {
	var order = models.Order{}
	order.UserID = user.ID
	order.ScooterID = scooter.ID

	querySQL := `INSERT INTO orders(user_id, scooter_id) VALUES ($1, $2) RETURNING id`
	err := ordb.db.QueryResultRow(context.Background(), querySQL, user.ID, scooter.ID).Scan(order.ID)
	if err != nil {
		return order, err
	}

	return order, nil
}

func (ordb *OrderRepoDb) SetOrderStart(order *models.Order, status models.ScooterStatusInRent) error {

	querySQL := `UPDATE orders(status_start_id) 
					SET status_start_id=$1
					WHERE id=$2`

	_, err := ordb.db.QueryResult(context.Background(), querySQL, status.ID, order.ID)

	return err
}