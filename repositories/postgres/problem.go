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

type ProblemRepoDB struct {
	userRepo    *UserRepoDB
	scooterRepo *ScooterRepoDB
	db          repositories.AnyDatabase
}

type SolutionRepoDB struct {
	db repositories.AnyDatabase
}

func NewProblemRepoDB(userRepo *UserRepoDB, scooterRepo *ScooterRepoDB, db repositories.AnyDatabase) *ProblemRepoDB {
	return &ProblemRepoDB{userRepo, scooterRepo, db}
}

func NewSolutionRepoDB(db repositories.AnyDatabase) *SolutionRepoDB {
	return &SolutionRepoDB{db}
}

func (probl *ProblemRepoDB) AddNewProblem(problem *models.Problem) error {
	querySQL := `INSERT INTO problems(user_id, type_Id, scooter_id, description, is_solved) 
		VALUES($1, $2, $3, $4, $5)
		RETURNING id, date_reported;`

	err := probl.db.QueryResultRow(context.Background(), querySQL,
		problem.User.ID, problem.Type.ID, problem.Scooter.ID, problem.Description, problem.IsSolved).
		Scan(&problem.ID, &problem.DateReported)

	return err
}

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

func (probl *ProblemRepoDB) MarkProblemAsSolved(problem *models.Problem) (models.Problem, error) {
	solvedWas := problem.IsSolved
	problem.IsSolved = true

	querySQL := `UPDATE problems SET is_solved = $1 WHERE id = $2`
	err := probl.db.QueryResultRow(context.Background(), querySQL, problem.IsSolved, problem.ID).Scan()

	if err != nil {
		problem.IsSolved = solvedWas
	}

	return *problem, err
}

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

func (probl *ProblemRepoDB) GetProblemTypeByID(typeID int) (models.ProblemType, error) {
	querySQL := `SELECT id, name FROM problem_types WHERE id = $1;`
	row := probl.db.QueryResultRow(context.Background(), querySQL, typeID)

	problemType := models.ProblemType{}
	err := row.Scan(&problemType.ID, &problemType.Name)

	return problemType, err
}

func (probl *ProblemRepoDB) GetProblemsByUserID(userID int) (*models.ProblemList, error) {
	return probl.getProblemsWithCondition(`user_id = $1`, userID)
}

func (probl *ProblemRepoDB) GetProblemsByTypeID(typeID int) (*models.ProblemList, error) {
	return probl.getProblemsWithCondition(`type_Id = $1`, typeID)
}

func (probl *ProblemRepoDB) GetProblemsByBeingSolved(solved bool) (*models.ProblemList, error) {
	return probl.getProblemsWithCondition(`is_solved = $1`, solved)
}

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

func (sol *SolutionRepoDB) AddProblemSolution(problemID int, solution *models.Solution) error {
	querySQL := `INSERT INTO solutions(problem_id, description) 
		VALUES($1, $2)
		RETURNING date_solved;`

	err := sol.db.QueryResultRow(context.Background(), querySQL, problemID, solution.Description).
		Scan(&solution.DateSolved)

	return err
}

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
