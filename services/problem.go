package services

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
	"time"
)

type ProblemService struct {
	repoProblem repositories.ProblemRepo
}

func NewProblemService(repoProblem repositories.ProblemRepo) *ProblemService {
	return &ProblemService{repoProblem}
}

func (problserv *ProblemService) AddNewProblem(problem *models.Problem) error {
	return problserv.repoProblem.AddNewProblem(problem)
}

func (problserv *ProblemService) GetProblemByID(problemID int) (models.Problem, error) {
	return problserv.repoProblem.GetProblemByID(problemID)
}

func (problserv *ProblemService) GetProblemTypeByID(typeID int) (models.ProblemType, error) {
	return problserv.repoProblem.GetProblemTypeByID(typeID)
}

func (problserv *ProblemService) GetProblemsByUserID(userID int) (*models.ProblemList, error) {
	return problserv.repoProblem.GetProblemsByUserID(userID)
}

func (problserv *ProblemService) GetProblemsByTypeID(typeID int) (*models.ProblemList, error) {
	return problserv.repoProblem.GetProblemsByTypeID(typeID)
}

func (problserv *ProblemService) GetProblemsByBeingSolved(solved bool) (*models.ProblemList, error) {
	return problserv.repoProblem.GetProblemsByBeingSolved(solved)
}

func (problserv *ProblemService) GetProblemsByTimePeriod(start, end time.Time) (*models.ProblemList, error) {
	return problserv.repoProblem.GetProblemsByTimePeriod(start, end)
}

func (problserv *ProblemService) AddProblemComplexFields(problem *models.Problem, typeID, scooterID, userID int) error {
	return problserv.repoProblem.AddProblemComplexFields(problem, typeID, scooterID, userID)
}