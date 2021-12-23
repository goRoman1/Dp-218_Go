package services

import (
	"Dp218Go/models"
	"Dp218Go/repositories/mock"
	"github.com/golang/mock/gomock"
	assert "github.com/stretchr/testify/require"
	"testing"
)

type problemUseCasesMock struct {
	repoProblem  *mock.MockProblemRepo
	repoSolution *mock.MockSolutionRepo
	problemUC    *ProblemService
}

type problemTestCase struct {
	name string
	test func(t *testing.T, mock *problemUseCasesMock)
}

func runProblemTestCases(t *testing.T, testCases []problemTestCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					tt.Error(err)
				}
			}()

			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			mock := newProblemUseCasesMock(ctrl)

			tc.test(tt, mock)
		})
	}
}

func newProblemUseCasesMock(ctrl *gomock.Controller) *problemUseCasesMock {
	repoProblem := mock.NewMockProblemRepo(ctrl)
	repoSolution := mock.NewMockSolutionRepo(ctrl)
	problemUC := NewProblemService(repoProblem, repoSolution)

	return &problemUseCasesMock{
		repoProblem:  repoProblem,
		repoSolution: repoSolution,
		problemUC:    problemUC,
	}
}

func Test_Problem_GetProblemByID(t *testing.T) {
	runProblemTestCases(t, []problemTestCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *problemUseCasesMock) {

				mock.repoProblem.EXPECT().GetProblemByID(1).
					Return(models.Problem{ID: 1}, nil).Times(1)

				result, err := mock.problemUC.GetProblemByID(1)
				assert.Equal(t, nil, err)
				assert.Equal(t, 1, result.ID)
			},
		},
	})
}

func Test_Problem_AddNewProblem(t *testing.T) {
	modelToReturn := &models.Problem{ID: 1, Description: "Test message"}

	runProblemTestCases(t, []problemTestCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *problemUseCasesMock) {
				mock.repoProblem.EXPECT().AddNewProblem(modelToReturn).
					Return(nil).Times(1)

				err := mock.problemUC.AddNewProblem(modelToReturn)
				assert.Equal(t, nil, err)
				assert.Equal(t, 1, modelToReturn.ID)
				assert.Contains(t, modelToReturn.Description, "Test")
			},
		},
	})
}