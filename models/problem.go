package models

import "time"

type ProblemType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Problem struct {
	ID           int         `json:"id"`
	User         User        `json:"user"`
	Type         ProblemType `json:"type"`
	Scooter      ScooterDTO  `json:"scooter"`
	DateReported time.Time   `json:"date_reported"`
	Description  string      `json:"description"`
	IsSolved     bool        `json:"is_solved"`
}

type ProblemList struct {
	Problems []Problem `json:"accounts"`
}
