package tests

import (
	"Dp218Go/models"
	"Dp218Go/services"
	"Dp218Go/tests/mocks"
	"errors"
	"github.com/golang/mock/gomock"
	"testing"
)

var ErrNoRows = errors.New("no rows in result set")

type usecasesUser struct {
	repoUser *mocks.MockUserRepo
	repoRole *mocks.MockRoleRepo
	userUC   *services.UserService
}

type testCase struct {
	name string
	test func(t *testing.T, mock *usecasesUser)
}

func runTestCases(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					tt.Error(err)
				}
			}()

			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			mock := NewUseCasesMock(ctrl)

			tc.test(tt, mock)
		})
	}
}

func NewUseCasesMock(ctrl *gomock.Controller) *usecasesUser {
	repoUser := mocks.NewMockUserRepo(ctrl)
	repoRole := mocks.NewMockRoleRepo(ctrl)
	userUC := services.NewUserService(repoUser, repoRole)

	return &usecasesUser{
		repoRole: repoRole,
		repoUser: repoUser,
		userUC:   userUC,
	}
}

func TestUseCases_Role_GetRoleByID(t *testing.T) {
	runTestCases(t, []testCase{
		{
			name: "correct",
			test: func(t *testing.T, mock *usecasesUser) {

				mock.repoRole.EXPECT().GetRoleByID(gomock.Eq(1)).
					Return(models.Role{}, nil).Times(1)

				_, err := mock.repoRole.GetRoleByID(1)
				if err != nil {
					t.Errorf("error in role repo -> GetRoleByID")
				}
			},
		},
		{
			name: "incorrect",
			test: func(t *testing.T, mock *usecasesUser) {

				mock.repoRole.EXPECT().GetRoleByID(gomock.Eq(0)).
					Return(models.Role{}, ErrNoRows).Times(1)

				_, err := mock.repoRole.GetRoleByID(0)
				if err == nil {
					t.Errorf("error expected but not received in role repo -> GetRoleByID")
				}
			},
		},
	})
}
