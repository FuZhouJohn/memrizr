package service

import (
	"context"
	"log"

	"github.com/FuZhouJohn/memrizr/account/model"
	"github.com/FuZhouJohn/memrizr/account/model/apperrors"
	"github.com/google/uuid"
)

type userService struct {
	UserRepository model.UserRepository
}

type USConfig struct {
	UserRepository model.UserRepository
}

func NewUserService(c *USConfig) model.UserService {
	return &userService{
		UserRepository: c.UserRepository,
	}
}

func (s *userService) Get(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	u, err := s.UserRepository.FindByID(ctx, uid)
	return u, err
}

func (s *userService) Signup(ctx context.Context, u *model.User) error {
	pw, err := hashPassword(u.Password)

	if err != nil {
		log.Printf("无法注册用户的电子邮件: %v\n", u.Email)
		return apperrors.NewInternal()
	}

	u.Password = pw

	if err := s.UserRepository.Create(ctx, u); err != nil {
		return err
	}

	return nil
}

func (s *userService) Signin(ctx context.Context, u *model.User) error {
	uFetched, err := s.UserRepository.FindByEmail(ctx, u.Email)
	if err != nil {
		return apperrors.NewAuthorization("用户名或密码错误")
	}

	match, err := comparePasswords(uFetched.Password, u.Password)
	if err != nil {
		return apperrors.NewInternal()
	}

	if !match {
		return apperrors.NewAuthorization("用户名或密码错误")
	}

	*u = *uFetched
	return nil
}
