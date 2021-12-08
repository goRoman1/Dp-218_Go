package postgres

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
	"context"
	"database/sql"
	"fmt"
)

type ScooterRepoDB struct {
	db repositories.AnyDatabase
}

func NewScooterRepoDB(db repositories.AnyDatabase) *ScooterRepoDB {
	return &ScooterRepoDB{db}
}

func (scdb *ScooterRepoDB) GetAllScooters () (*models.ScooterList, error) {
	scooterList := &models.ScooterList{}

	querySQL := `SELECT s.id, sm.max_weight, sm.model_name, ss.battery_remain, ss.latitude, ss.longitude 
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

	for rows.Next() {
		var scooter models.Scooter
		err := rows.Scan(&scooter.ID, &scooter.MaxWeight, &scooter.ScooterModel, &scooter.BatteryRemain,
			&scooter.Latitude, &scooter.Longitude)
		if err != nil {
			return scooterList, err
		}
		scooterList.Scooters = append(scooterList.Scooters, scooter)
		fmt.Println(scooter)
	}
	fmt.Println(scooterList)
	return scooterList, nil
}

func (scdb *ScooterRepoDB) GetScooterById(scooterId int) (models.Scooter, error) {
	scooter := models.Scooter{}
	querySQL := `SELECT s.id, sm.max_weight, sm.model_name, ss.battery_remain, ss.latitude, 
ss.longitude 
					FROM scooters as s 
					JOIN scooter_models as sm 
					ON s.model_id=sm.id 
					JOIN scooter_statuses as ss 
					ON s.id=ss.scooter_id 
					WHERE s.id=$1`

	row := scdb.db.QueryResultRow(context.Background(), querySQL, scooterId)
	switch err := row.Scan(&scooter.ID, &scooter.MaxWeight, &scooter.ScooterModel, &scooter.BatteryRemain,
		&scooter.Latitude, &scooter.Longitude); err {
	case sql.ErrNoRows:
		return scooter, err
	default:
		return scooter, err
	}
}

func (scdb *ScooterRepoDB) GetScooterStatus(scooterID int) (models.ScooterStatus, error) {
	var scooterStatus = models.ScooterStatus{}
	scooter, err := scdb.GetScooterById(scooterID)
	if err!=nil {
		fmt.Println(err)
		return models.ScooterStatus{}, err
	}
	scooterStatus.Scooter = scooter

	querySQL := `SELECT battery_remain, latitude, longitude 
				FROM scooter_statuses
				WHERE scooter_id=$1`

	row := scdb.db.QueryResultRow(context.Background(),querySQL, scooterID)
	switch err = row.Scan(&scooterStatus.BatteryRemain,
		&scooterStatus.Location.Latitude, &scooterStatus.Location.Longitude); err {
	case sql.ErrNoRows:
		return scooterStatus, err
	default:
		return scooterStatus, err
	}
}


//func(scdb *ScooterRepoDB) SendAtStart(uID, sID int) (error, int) {
//	scooter, err := scdb.GetScooterById(sID)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	coordinate := models.Coordinate{Latitude: scooter.Latitude, Longitude: scooter.Longitude}
//
//	var tripId int
//	querySQL := `INSERT INTO scooter_statuses_in_rent(user_id, scooter_id, date_time)
//					VALUES ($1, $2, now())
//					RETURNING id`
//	err = scdb.db.QueryResultRow(context.Background(), querySQL, uID, sID).Scan(&tripId)
//	if err != nil {
//		return err, 0
//	}
//
//	querySQL = `INSERT INTO locations(latitude, longitude, label)
//					VALUES($1, $2, $3)
//					RETURNING id`
//	_, err = scdb.db.QueryResult(context.Background(), querySQL, coordinate.Latitude, coordinate.Longitude,
//		string(rune(tripId)))
//	if err != nil {
//		return err, 0
//	}
//
//	return nil, tripId
//}

//func (scdb *ScooterRepoDB) SendAtEnd(tripId int, client *Client) error {
//	querySQL := `INSERT INTO locations(latitude, longitude, label)
//					VALUES($1, $2, $3)`
//	_, err := scdb.db.QueryResult(context.Background(), querySQL, client.Latitude, client.Longitude, string(rune(tripId)))
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	return nil
//}