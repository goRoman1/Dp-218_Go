package repositories

import (
	model "Dp218Go/domain/dto"
	"Dp218Go/pkg/postgres"
	"context"
	"database/sql"
)

type ScooterRepoDb struct {
	*postgres.Postgres
}

func NewSc(pg *postgres.Postgres) *ScooterRepoDb {
	return &ScooterRepoDb{pg}
}

func (sc *ScooterRepoDb) GetAllScooters () (*model.ScooterList, error) {
	scooterList := &model.ScooterList{}
	sc.QuerySQL = "SELECT s.id, sm.maxweight, sm.modelname, ss.locationid, ss.batteryremain, ss.canberent, l.lattitude, l.longtitude  FROM scooters as s JOIN scootermodels as sm ON s.modelid=sm.id JOIN scooterstatuses as ss ON s.id=ss.scooterid JOIN locations as l ON ss.locationid=l.id ORDER BY s.id"
	rows, err := sc.QueryResult(context.Background())
	if err != nil {
		return scooterList, err
	}
	for rows.Next() {
		var scooter model.Scooter
		err := rows.Scan(&scooter.Id, &scooter.MaxWeight, &scooter.ScooterModel,&scooter.LocationId, &scooter.BatteryRemain, &scooter.CanBeRent, &scooter.Lattitude, &scooter.Longtitude)
		if err != nil {
			return scooterList, err
		}
		scooterList.Scooters = append(scooterList.Scooters, scooter)
	}
	return scooterList, nil
}

func (sc *ScooterRepoDb) GetScooterById(scooterId int) (model.Scooter, error) {
	scooter := model.Scooter{}
	sc.QuerySQL = "SELECT s.id, sm.maxweight, sm.modelname, ss.locationid, ss.batteryremain, ss.canberent, l.lattitude, l.longtitude  FROM scooters as s JOIN scootermodels as sm ON s.modelid=sm.id JOIN scooterstatuses as ss ON s.id=ss.scooterid JOIN locations as l ON ss.locationid=l.id WHERE s.id=$1"
	row := sc.QueryResultRow(context.Background(), scooterId)
	switch err := row.Scan(&scooter.Id, &scooter.MaxWeight, &scooter.ScooterModel,
		&scooter.LocationId, &scooter.BatteryRemain, &scooter.CanBeRent, &scooter.Lattitude, &scooter.Longtitude); err {
	case sql.ErrNoRows:
		return scooter, ErrNoMatch
	default:
		return scooter, err
	}
}

//func (sc *ScooterRepoDb) SendPosition(scooter model.Scooter) {
//	sc.QuerySQL = "UPDATE locations SET lattitude=$1, longtitude=$2 WHERE id=$3 RETURNING id"
//	_, err := sc.QueryResult(context.Background(), scooter.Lattitude, scooter.Longtitude,scooter.LocationId )
//	if err!=nil {
//		fmt.Println(err)
//	}
//}
