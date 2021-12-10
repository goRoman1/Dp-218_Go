package models

import "time"

type ScooterDTO struct {
	ID            int     `json:"scooter_id"`
	ScooterModel  string  `json:"scooter_model"`
	MaxWeight     float64 `json:"max_weight"`
	BatteryRemain float64 `json:"battery_remain"`
	CanBeRent     bool    `json:"can_be_rent"`
}

type ScooterListDTO struct {
	Scooters []ScooterDTO `json:"scooters"`
}

type ScooterStatus struct {
	Scooter       ScooterDTO `json:"scooter"`
	Location      Coordinate `json:"location"`
	BatteryRemain int        `json:"battery_remain"`
	StationID     int        `json:"station_id"`
}

type ScooterStatusInRent struct {
	ID        int        `json:"id"`
	StationID int        `json:"station_id"`
	DateTime  time.Time  `json:"date_time"`
	Location  Coordinate `json:"location"`
}
