package repositories

import (
	"Dp218Go/models"
	"time"
)

type ProblemRepo interface {
	AddNewProblem(problem *models.Problem) error
	GetProblemByID(problemID int) (models.Problem, error)
	GetProblemTypeByID(typeID int) (models.ProblemType, error)
	GetProblemsByUserID(userID int) (*models.ProblemList, error)
	GetProblemsByTypeID(typeID int) (*models.ProblemList, error)
	GetProblemsByBeingSolved(solved bool) (*models.ProblemList, error)
	GetProblemsByTimePeriod(start, end time.Time) (*models.ProblemList, error)
	AddProblemComplexFields(problem *models.Problem, typeID, scooterID, userID int) error
}
