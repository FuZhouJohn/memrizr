package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/FuZhouJohn/memrizr/account/model"
	"github.com/FuZhouJohn/memrizr/account/model/apperrors"
	"github.com/FuZhouJohn/memrizr/account/model/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUserResp := &model.User{
			UID:   uid,
			Email: "bob@bob.com",
			Name:  "Bobby Bobson",
		}

		mockUserRespostiory := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRespostiory,
		})
		mockUserRespostiory.On("FindByID", mock.Anything, uid).Return(mockUserResp, nil)

		ctx := context.TODO()
		u, err := us.Get(ctx, uid)

		assert.NoError(t, err)
		assert.Equal(t, u, mockUserResp)
		mockUserRespostiory.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUserRespository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRespository,
		})

		mockUserRespository.On("FindByID", mock.Anything, uid).Return(nil, fmt.Errorf("Some error down the call chhain"))

		ctx := context.TODO()
		u, err := us.Get(ctx, uid)

		assert.Nil(t, u)
		assert.Error(t, err)
		mockUserRespository.AssertExpectations(t)
	})
}

func TestSignup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUser := &model.User{
			Email:    "hello@world.com",
			Password: "testpassword",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})
		mockUserRepository.On("Create", mock.AnythingOfType("*context.emptyCtx"), mockUser).
			Run(func(args mock.Arguments) {
				userArg := args.Get(1).(*model.User)
				userArg.UID = uid
			}).Return(nil)

		ctx := context.TODO()
		err := us.Signup(ctx, mockUser)
		assert.NoError(t, err)

		assert.Equal(t, uid, mockUser.UID)

		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUser := &model.User{
			Email:    "hello@world.com",
			Password: "testpassword",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})

		mockErr := apperrors.NewConflict("email", mockUser.Email)

		mockUserRepository.
			On("Create", mock.AnythingOfType("*context.emptyCtx"), mockUser).
			Return(mockErr)

		ctx := context.TODO()
		err := us.Signup(ctx, mockUser)

		assert.EqualError(t, err, mockErr.Error())

		mockUserRepository.AssertExpectations(t)
	})
}
