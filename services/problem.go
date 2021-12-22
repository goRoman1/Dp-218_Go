package services

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
	"time"
)

// ProblemService - structure for implementing user problem service
type ProblemService struct {
	repoProblem  repositories.ProblemRepo
	repoSolution repositories.SolutionRepo
}

// NewProblemService - initialization of ProblemService
func NewProblemService(repoProblem repositories.ProblemRepo, repoSolution repositories.SolutionRepo) *ProblemService {
	return &ProblemService{repoProblem, repoSolution}
}

// AddNewProblem - add new user problem record
func (problserv *ProblemService) AddNewProblem(problem *models.Problem) error {
	return problserv.repoProblem.AddNewProblem(problem)
}

// GetProblemByID - get problem information by its ID
func (problserv *ProblemService) GetProblemByID(problemID int) (models.Problem, error) {
	return problserv.repoProblem.GetProblemByID(problemID)
}

// MarkProblemAsSolved - update problem record to make problem solved
func (problserv *ProblemService) MarkProblemAsSolved(problem *models.Problem) (models.Problem, error) {
	return problserv.repoProblem.MarkProblemAsSolved(problem)
}

// GetProblemTypeByID - get problem type record by its ID
func (problserv *ProblemService) GetProblemTypeByID(typeID int) (models.ProblemType, error) {
	return problserv.repoProblem.GetProblemTypeByID(typeID)
}

func (problserv *ProblemService) GetAllProblemTypes() ([]models.ProblemType, error) {
	return problserv.repoProblem.GetAllProblemTypes()
}

// GetProblemsByUserID - get problem list for given user (by user ID)
func (problserv *ProblemService) GetProblemsByUserID(userID int) (*models.ProblemList, error) {
	return problserv.repoProblem.GetProblemsByUserID(userID)
}

// GetProblemsByTypeID - get problem list by given problem type ID
func (problserv *ProblemService) GetProblemsByTypeID(typeID int) (*models.ProblemList, error) {
	return problserv.repoProblem.GetProblemsByTypeID(typeID)
}

// GetProblemsByBeingSolved - get problem list by is_solved field value
func (problserv *ProblemService) GetProblemsByBeingSolved(solved bool) (*models.ProblemList, error) {
	return problserv.repoProblem.GetProblemsByBeingSolved(solved)
}

// GetProblemsByTimePeriod - get problem list from time start to time end
func (problserv *ProblemService) GetProblemsByTimePeriod(start, end time.Time) (*models.ProblemList, error) {
	return problserv.repoProblem.GetProblemsByTimePeriod(start, end)
}

// AddProblemComplexFields - fulfill problem model with problem type, scooter, user (by their IDs)
func (problserv *ProblemService) AddProblemComplexFields(problem *models.Problem, typeID, scooterID, userID int) error {
	return problserv.repoProblem.AddProblemComplexFields(problem, typeID, scooterID, userID)
}

// AddProblemSolution - make solution record for given problem (by ID)
func (problserv *ProblemService) AddProblemSolution(problemID int, solution *models.Solution) error {
	err := problserv.repoSolution.AddProblemSolution(problemID, solution)
	if err != nil {
		return err
	}
	problem, err := problserv.repoProblem.GetProblemByID(problemID)
	if err != nil {
		return err
	}
	_, err = problserv.repoProblem.MarkProblemAsSolved(&problem)
	return err
}

// GetSolutionByProblem - get solution for given problem
func (problserv *ProblemService) GetSolutionByProblem(problem models.Problem) (models.Solution, error) {
	return problserv.repoSolution.GetSolutionByProblem(problem)
}

// GetSolutionsByProblems - get solutions for given problems list
func (problserv *ProblemService) GetSolutionsByProblems(problems models.ProblemList) (map[models.Problem]models.Solution, error) {
	return problserv.repoSolution.GetSolutionsByProblems(problems)
}
