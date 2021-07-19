package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/FuZhouJohn/memrizr/account/model"
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

		mockUserRespostiory := new(mocks.MockUserRespository)
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

		mockUserRespository := new(mocks.MockUserRespository)
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
