package postgres

import (
	"Dp218Go/models"
	"context"
)

func (pg *Postgres) GetAllStations() (*models.StationList, error) {
	list := &models.StationList{}

	querySQL := `SELECT * FROM scooter_stations ORDER BY id DESC;`
	rows, err := pg.QueryResult(context.Background(), querySQL)
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var station models.Station
		err := rows.Scan(&station.ID, &station.LocationID, &station.Name, &station.IsActive)
		if err != nil {
			return list, err
		}

		list.Station = append(list.Station, station)
	}
	return list, nil
}

func (pg *Postgres) AddStation(station *models.Station) error {
	var id int
	querySQL := `INSERT INTO scooter_stations(id, location_id, name, is_active) 
		VALUES($1, $2, $3, $4)
		RETURNING id;`
	err := pg.QueryResultRow(context.Background(), querySQL, station.ID, station.LocationID, station.Name, station.IsActive).Scan(&id)
	if err != nil {
		return err
	}
	station.ID = id
	return nil
}

func (pg *Postgres) GetStationById(stationId int) (models.Station, error) {
	station := models.Station{}

	querySQL := `SELECT * FROM scooter_stations WHERE id = $1;`
	row := pg.QueryResultRow(context.Background(), querySQL, stationId)
	err := row.Scan(&station.ID, &station.LocationID, &station.Name, &station.IsActive)

	return station, err
}

func (pg *Postgres) DeleteStation(stationId int) error {
	querySQL := `DELETE FROM scooter_stations WHERE id = $1;`
	_, err := pg.QueryExec(context.Background(), querySQL, stationId)
	return err
}