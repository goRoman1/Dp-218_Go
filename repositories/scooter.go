package repositories

import (
	model "Dp218Go/domain/dto"
	"Dp218Go/pkg/postgres"
	"context"
	"database/sql"
	"fmt"
)

type ScooterRepoDb struct {
	*postgres.Postgres
}

func NewSc(pg *postgres.Postgres) *ScooterRepoDb {
	return &ScooterRepoDb{pg}
}

func (sc *ScooterRepoDb) GetAllScooters () (*model.ScooterList, error) {
	scooterList := &model.ScooterList{}
	sc.QuerySQL = `SELECT s.id, sm.max_weight, sm.model_name, ss.battery_remain, ss.can_be_rent, s.latitude, s.longitude FROM scooters as s JOIN scooter_models as sm ON s.model_id=sm.id JOIN scooter_statuses as ss ON s.id=ss.scooter_id ORDER BY s.id`
	rows, err := sc.QueryResult(context.Background())
	if err != nil {
		return scooterList, err
	}
	for rows.Next() {
		var scooter model.Scooter
		err := rows.Scan(&scooter.Id, &scooter.MaxWeight, &scooter.ScooterModel, &scooter.BatteryRemain, &scooter.CanBeRent, &scooter.Latitude, &scooter.Longitude)
		if err != nil {
			return scooterList, err
		}
		scooterList.Scooters = append(scooterList.Scooters, scooter)
	}
	return scooterList, nil
}

func (sc *ScooterRepoDb) GetScooterById(scooterId int) (model.Scooter, error) {
	scooter := model.Scooter{}
	sc.QuerySQL = "SELECT s.id, sm.max_weight, sm.model_name, ss.location_id, ss.battery_remain, ss.can_be_rent, l.latitude, l.longitude  FROM scooters as s JOIN scooter_models as sm ON s.model_id=sm.id JOIN scooter_statuses as ss ON s.id=ss.scooter_id JOIN locations as l ON ss.location_id=l.id WHERE s.id=$1"
	row := sc.QueryResultRow(context.Background(), scooterId)
	switch err := row.Scan(&scooter.Id, &scooter.MaxWeight, &scooter.ScooterModel,
		&scooter.LocationId, &scooter.BatteryRemain, &scooter.CanBeRent, &scooter.Latitude, &scooter.Longitude); err {
	case sql.ErrNoRows:
		return scooter, ErrNoMatch
	default:
		return scooter, err
	}
}

func (sc *ScooterRepoDb) SendPosition(scooter model.Scooter) {
	sc.QuerySQL = "UPDATE locations SET latitude=$1, longitude=$2 WHERE id=$3 RETURNING id"
	_, err := sc.QueryResult(context.Background(), scooter.Latitude, scooter.Longitude,scooter.LocationId )
	if err!=nil {
		fmt.Println(err)
	}
}
