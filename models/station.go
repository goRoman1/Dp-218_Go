package models


type Station struct {
	ID          int    `json:"id"`
	LocationID int `json:"location_id"`
	Name string `json:"name"`
	IsActive bool `json:"is_active"`
}

type StationList struct {
	Station []Station `json:"station"`
}