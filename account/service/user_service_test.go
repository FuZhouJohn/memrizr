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
	t.Run("成功", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUser := &model.User{
			Email:    "hello@world1.com",
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
			Email:    "hello@world2.com",
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

func TestSignin(t *testing.T) {
	email := "hello@world2.com"
	vaildPW := "zhuangjinan"
	hashedVaildPW, _ := hashPassword(vaildPW)
	invalidPW := "zhuangjibei"

	mockeUserRepository := new(mocks.MockUserRepository)
	us := NewUserService(&USConfig{
		UserRepository: mockeUserRepository,
	})

	t.Run("成功", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUser := &model.User{
			Email:    email,
			Password: vaildPW,
		}

		mockUserResp := &model.User{
			UID:      uid,
			Email:    email,
			Password: hashedVaildPW,
		}

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			email,
		}

		mockeUserRepository.On("FindByEmail", mockArgs...).Return(mockUserResp, nil)

		ctx := context.TODO()
		err := us.Signin(ctx, mockUser)

		assert.NoError(t, err)
		mockeUserRepository.AssertCalled(t, "FindByEmail", mockArgs...)
	})

	t.Run("无效的电子邮件/密码组合", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUser := &model.User{
			Email:    email,
			Password: invalidPW,
		}

		mockUserResp := &model.User{
			UID:      uid,
			Email:    email,
			Password: hashedVaildPW,
		}

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			email,
		}

		mockeUserRepository.On("FindByEmail", mockArgs...).Return(mockUserResp, nil)

		ctx := context.TODO()
		err := us.Signin(ctx, mockUser)

		assert.Error(t, err)
		assert.EqualError(t, err, "用户名或密码错误")
		mockeUserRepository.AssertCalled(t, "FindByEmail", mockArgs...)
	})
}
