package postgres

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
	"context"
	"time"
)

type ProblemRepoDB struct {
	userRepo    *UserRepoDB
	scooterRepo *ScooterRepoDB
	db          repositories.AnyDatabase
}

func NewProblemRepoDB(userRepo *UserRepoDB, scooterRepo *ScooterRepoDB, db repositories.AnyDatabase) *ProblemRepoDB {
	return &ProblemRepoDB{userRepo, scooterRepo, db}
}

func (probl *ProblemRepoDB) AddNewProblem(problem *models.Problem) error {
	querySQL := `INSERT INTO problems(user_id, type_Id, scooter_id, description, is_solved) 
		VALUES($1, $2, $3, $4, $5)
		RETURNING id, date_reported;`

	err := probl.db.QueryResultRow(context.Background(), querySQL,
		problem.User.ID, problem.Type.ID, problem.Scooter.ID, problem.Description, problem.IsSolved).
		Scan(&problem.ID, &problem.DateReported)
	if err != nil {
		return err
	}

	return nil
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
