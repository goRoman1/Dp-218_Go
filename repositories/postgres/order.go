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

func (ordb *OrderRepoDb) CreateOrder(user models.User, scooterID, startID, endID int, distance float64) (models.Order,
	error) {
	var order = models.Order{}
	order.UserID = user.ID
	order.ScooterID = scooterID
	order.StatusStartID = startID
	order.StatusEndID = endID
	order.Distance = distance

	querySQL := `INSERT INTO orders(user_id, scooter_id, status_start_id, status_end_id, distance) 
					VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := ordb.db.QueryResultRow(context.Background(), querySQL, user.ID, scooterID, startID, endID,
		distance).Scan(&order.ID)
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

func (ordb *OrderRepoDb) GetAllOrders() (*models.OrderList, error) {
	orderList := &models.OrderList{}

	querySQL := `SELECT * from orders`
	rows, err := ordb.db.QueryResult(context.Background(), querySQL)
	if err!=nil {
		return orderList, err
	}
	defer rows.Close()
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.ScooterID, &order.StatusStartID, &order.StatusEndID,
			&order.Distance, &order.Amount)
		if err != nil {
			return orderList, err
		}
		orderList.Orders = append(orderList.Orders, order)
	}
	return orderList, nil
}