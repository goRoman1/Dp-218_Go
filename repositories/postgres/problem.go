package postgres

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ProblemRepoDB - struct for user problems repository
type ProblemRepoDB struct {
	userRepo    *UserRepoDB
	scooterRepo *ScooterRepoDB
	db          repositories.AnyDatabase
}

// SolutionRepoDB - struct for solutions repository
type SolutionRepoDB struct {
	db repositories.AnyDatabase
}

// NewProblemRepoDB - initialize problem repo from user & scooter repos
func NewProblemRepoDB(userRepo *UserRepoDB, scooterRepo *ScooterRepoDB, db repositories.AnyDatabase) *ProblemRepoDB {
	return &ProblemRepoDB{userRepo, scooterRepo, db}
}

// NewSolutionRepoDB - initialize solution repo
func NewSolutionRepoDB(db repositories.AnyDatabase) *SolutionRepoDB {
	return &SolutionRepoDB{db}
}

// AddNewProblem - create new user problem record in the DB based on problem entity
func (probl *ProblemRepoDB) AddNewProblem(problem *models.Problem) error {
	querySQL := `INSERT INTO problems(user_id, type_Id, scooter_id, description, is_solved) 
		VALUES($1, $2, $3, $4, $5)
		RETURNING id, date_reported;`

	err := probl.db.QueryResultRow(context.Background(), querySQL,
		problem.User.ID, problem.Type.ID, problem.Scooter.ID, problem.Description, problem.IsSolved).
		Scan(&problem.ID, &problem.DateReported)

	return err
}

// GetProblemByID - get user problem record from the DB by given problem ID
func (probl *ProblemRepoDB) GetProblemByID(problemID int) (models.Problem, error) {
	problem := models.Problem{}

	querySQL := `SELECT 
		id, user_id, type_Id, scooter_id, date_reported, description, is_solved
		FROM problems 
		WHERE id = $1;`
	row := probl.db.QueryResultRow(context.Background(), querySQL, problemID)

	var userID, typeID, scooterID int
	err := row.Scan(&problem.ID, &userID, &typeID, &scooterID,
		&problem.DateReported, &problem.Description, &problem.IsSolved)
	if err != nil {
		return models.Problem{}, err
	}

	err = probl.AddProblemComplexFields(&problem, typeID, scooterID, userID)
	if err != nil {
		return models.Problem{}, err
	}

	return problem, err
}

// MarkProblemAsSolved - change DB record of given problem to mark it as solved
func (probl *ProblemRepoDB) MarkProblemAsSolved(problem *models.Problem) (models.Problem, error) {
	solvedWas := problem.IsSolved
	problem.IsSolved = true

	querySQL := `UPDATE problems SET is_solved = $1 WHERE id = $2 RETURNING is_solved;`
	err := probl.db.QueryResultRow(context.Background(), querySQL, problem.IsSolved, problem.ID).Scan(&problem.IsSolved)

	if err != nil {
		problem.IsSolved = solvedWas
	}

	return *problem, err
}

// AddProblemComplexFields - add data about problem type, scooter, user to given problem entity by given IDs
func (probl *ProblemRepoDB) AddProblemComplexFields(problem *models.Problem, typeID, scooterID, userID int) error {
	var err error
	if userID != 0 {
		problem.User, err = probl.userRepo.GetUserByID(userID)
		if err != nil {
			return err
		}
	}
	if scooterID != 0 {
		problem.Scooter, err = probl.scooterRepo.GetScooterById(scooterID)
		if err != nil {
			return err
		}
	}
	if typeID != 0 {
		problem.Type, err = probl.GetProblemTypeByID(typeID)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetProblemTypeByID - get problem type by given ID from the DB
func (probl *ProblemRepoDB) GetProblemTypeByID(typeID int) (models.ProblemType, error) {
	querySQL := `SELECT id, name FROM problem_types WHERE id = $1;`
	row := probl.db.QueryResultRow(context.Background(), querySQL, typeID)

	problemType := models.ProblemType{}
	err := row.Scan(&problemType.ID, &problemType.Name)

	return problemType, err
}

// GetAllProblemTypes - get all available problem types from the DB
func (probl *ProblemRepoDB) GetAllProblemTypes() ([]models.ProblemType, error) {
	var result []models.ProblemType
	querySQL := `SELECT 
		id, name 
		FROM problem_types`

	rows, err := probl.db.QueryResult(context.Background(), querySQL)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var problemType models.ProblemType
		err := rows.Scan(&problemType.ID, &problemType.Name)
		if err != nil {
			return result, err
		}
		result = append(result, problemType)
	}

	return result, nil
}

// GetProblemsByUserID - get list of problems from the DB by user ID
func (probl *ProblemRepoDB) GetProblemsByUserID(userID int) (*models.ProblemList, error) {
	return probl.getProblemsWithCondition(`user_id = $1`, userID)
}

// GetProblemsByTypeID - get list of problems from the DB by type ID
func (probl *ProblemRepoDB) GetProblemsByTypeID(typeID int) (*models.ProblemList, error) {
	return probl.getProblemsWithCondition(`type_Id = $1`, typeID)
}

// GetProblemsByBeingSolved - get list of problems from the DB by is_solved field
func (probl *ProblemRepoDB) GetProblemsByBeingSolved(solved bool) (*models.ProblemList, error) {
	return probl.getProblemsWithCondition(`is_solved = $1`, solved)
}

// GetProblemsByTimePeriod - get list of problems from the DB from start to end time
func (probl *ProblemRepoDB) GetProblemsByTimePeriod(start, end time.Time) (*models.ProblemList, error) {
	return probl.getProblemsWithCondition(`date_reported >= $1 AND date_reported <= $2`, start, end)
}

func (probl *ProblemRepoDB) getProblemsWithCondition(condition string, params ...interface{}) (*models.ProblemList, error) {
	list := &models.ProblemList{}

	querySQL := `SELECT 
		id, user_id, type_Id, scooter_id, date_reported, description, is_solved 
		FROM problems 
		WHERE ` + condition + `
		ORDER BY date_reported;`
	rows, err := probl.db.QueryResult(context.Background(), querySQL, params...)
	if err != nil {
		return list, err
	}
	defer rows.Close()

	type additionalProblemData struct {
		typeID    int
		scooterID int
		userID    int
	}
	var problemAdditionalData = make(map[models.Problem]additionalProblemData)

	for rows.Next() {
		var problem models.Problem
		var typeID, scooterID, userID int
		err := rows.Scan(&problem.ID, &userID, &typeID, &scooterID,
			&problem.DateReported, &problem.Description, &problem.IsSolved)
		if err != nil {
			return list, err
		}

		problemAdditionalData[problem] = additionalProblemData{
			typeID:    typeID,
			scooterID: scooterID,
			userID:    userID,
		}
	}

	for key, value := range problemAdditionalData {
		err = probl.AddProblemComplexFields(&key, value.typeID, value.scooterID, value.userID)
		if err != nil {
			return list, err
		}

		list.Problems = append(list.Problems, key)
	}

	return list, nil
}

// AddProblemSolution - create new DB record for problem solution in the DB by problem ID & solution info
func (sol *SolutionRepoDB) AddProblemSolution(problemID int, solution *models.Solution) error {
	querySQL := `INSERT INTO solutions(problem_id, description) 
		VALUES($1, $2)
		RETURNING date_solved;`

	err := sol.db.QueryResultRow(context.Background(), querySQL, problemID, solution.Description).
		Scan(&solution.DateSolved)

	return err
}

// GetSolutionByProblem - get solution info for given problem from the DB
func (sol *SolutionRepoDB) GetSolutionByProblem(problem models.Problem) (models.Solution, error) {
	solution := models.Solution{}
	solution.Problem = problem

	querySQL := `SELECT 
		description, date_solved
		FROM solutions 
		WHERE problem_id = $1;`
	err := sol.db.QueryResultRow(context.Background(), querySQL, problem.ID).
		Scan(&solution.Description, &solution.DateSolved)

	return solution, err
}

// GetSolutionsByProblems - get solutions for given problems list from the DB
func (sol *SolutionRepoDB) GetSolutionsByProblems(problems models.ProblemList) (map[models.Problem]models.Solution, error) {
	result := make(map[models.Problem]models.Solution, len(problems.Problems))
	if len(result) == 0 {
		return result, fmt.Errorf("no problems to get solution to")
	}

	querySQL := `SELECT 
		problem_id, description, date_solved
		FROM solutions 
		WHERE problem_id IN `
	params := make([]int, len(problems.Problems))
	conditions := make([]string, len(problems.Problems))
	for i := 0; i < len(problems.Problems); i++ {
		params[i] = problems.Problems[i].ID
		conditions[i] = "$" + strconv.Itoa(i+1)
	}
	querySQL += "(" + strings.Join(conditions, ",") + ")"
	querySQL += ";"

	rows, err := sol.db.QueryResult(context.Background(), querySQL, params)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var solution models.Solution
		var problemID int

		err := rows.Scan(&problemID, &solution.Description, &solution.DateSolved)
		if err != nil {
			return result, err
		}
		solution.Problem, err = findProblemByID(problems, problemID)
		if err != nil {
			return result, err
		}

		result[solution.Problem] = solution
	}

	return result, nil
}

func findProblemByID(problems models.ProblemList, problemID int) (models.Problem, error) {
	for _, v := range problems.Problems {
		if v.ID == problemID {
			return v, nil
		}
	}
	return models.Problem{}, fmt.Errorf("no such problem found in list")
}
