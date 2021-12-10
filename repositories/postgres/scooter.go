package postgres

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
	"context"
	"fmt"
)

type ScooterRepoDB struct {
	db repositories.AnyDatabase
}

func NewScooterRepoDB(db repositories.AnyDatabase) *ScooterRepoDB {
	return &ScooterRepoDB{db}
}

func (scdb *ScooterRepoDB) GetAllScooters() (*models.ScooterListDTO, error) {
	scooterList := &models.ScooterListDTO{}

	querySQL := `SELECT s.id, sm.max_weight, sm.model_name, ss.battery_remain, ss.can_be_rent
					FROM scooters as s 
					JOIN scooter_models as sm 
					ON s.model_id=sm.id 
					JOIN scooter_statuses as ss 
					ON s.id=ss.scooter_id 
					ORDER BY s.id`

	rows, err := scdb.db.QueryResult(context.Background(), querySQL)
	if err != nil {
		return scooterList, err
	}
	defer rows.Close()

	for rows.Next() {
		var scooter models.ScooterDTO
		err := rows.Scan(&scooter.ID, &scooter.MaxWeight, &scooter.ScooterModel, &scooter.BatteryRemain, &scooter.CanBeRent)
		if err != nil {
			return scooterList, err
		}
		scooterList.Scooters = append(scooterList.Scooters, scooter)
	}
	return scooterList, nil
}

func (scdb *ScooterRepoDB) GetScooterById(scooterId int) (models.ScooterDTO, error) {
	scooter := models.ScooterDTO{}
	querySQL := `SELECT s.id, sm.max_weight, sm.model_name, ss.battery_remain, ss.can_be_rent
					FROM scooters as s 
					JOIN scooter_models as sm 
					ON s.model_id=sm.id 
					JOIN scooter_statuses as ss 
					ON s.id=ss.scooter_id 
					WHERE s.id=$1`

	row := scdb.db.QueryResultRow(context.Background(), querySQL, scooterId)
	err := row.Scan(&scooter.ID, &scooter.MaxWeight, &scooter.ScooterModel, &scooter.BatteryRemain, &scooter.CanBeRent)
	if err !=nil {
		return scooter, err
	}

	return scooter, nil
}

func (scdb *ScooterRepoDB) GetScooterStatus(scooterID int) (models.ScooterStatus, error) {
	var scooterStatus = models.ScooterStatus{}
	scooter, err := scdb.GetScooterById(scooterID)
	if err != nil {
		fmt.Println(err)
		return models.ScooterStatus{}, err
	}
	scooterStatus.Scooter = scooter

	querySQL := `SELECT battery_remain, latitude, longitude 
					FROM scooter_statuses
					WHERE scooter_id=$1`

	row := scdb.db.QueryResultRow(context.Background(), querySQL, scooterID)
	err = row.Scan(&scooterStatus.BatteryRemain,
		&scooterStatus.Location.Latitude, &scooterStatus.Location.Longitude)
	if err != nil {
		return scooterStatus, err
	}

	return scooterStatus, nil
}

func (scdb *ScooterRepoDB) CreateScooterStatusInRent(scooterID int) (models.ScooterStatusInRent, error) {
	var scooterStatusInRent models.ScooterStatusInRent
	scooterStatus, err := scdb.GetScooterStatus(scooterID)
	if err != nil {
		fmt.Println(err)
		return scooterStatusInRent, err
	}

	scooterStatusInRent.Location = scooterStatus.Location

	querySQL := `INSERT INTO scooter_statuses_in_rent(date_time, latitude, longitude) 
					VALUES(now(), $1, $2) RETURNING id, date_time`

	err = scdb.db.QueryResultRow(context.Background(), querySQL, scooterStatus.Location.Latitude,
		scooterStatus.Location.Longitude).Scan(&scooterStatusInRent.ID, &scooterStatusInRent.DateTime)
	if err != nil {
		fmt.Println(err)
		return scooterStatusInRent, err
	}

	return scooterStatusInRent, nil

}

func (scdb *ScooterRepoDB) SendCurrentPosition(id int, lat, lon float64) error {
	querySQL := `UPDATE scooter_statuses 
					SET latitude=$1, longitude=$2
					WHERE scooter_id=$3`

	_, err := scdb.db.QueryResult(context.Background(), querySQL, lat, lon, id)
	return err
}
