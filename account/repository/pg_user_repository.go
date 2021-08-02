package repository

import (
	"context"
	"log"

	"github.com/FuZhouJohn/memrizr/account/model"
	"github.com/FuZhouJohn/memrizr/account/model/apperrors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type pgUserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) model.UserRepository {
	return &pgUserRepository{
		DB: db,
	}
}

func (r *pgUserRepository) Create(ctx context.Context, u *model.User) error {
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *"

	if err := r.DB.GetContext(ctx, u, query, u.Email, u.Password); err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			log.Printf("无法使用该邮箱创建用户：%v。原因是：%v\n", u.Email, err.Code.Name())
			return apperrors.NewConflict("email", u.Email)
		}

		log.Printf("无法使用该邮箱创建用户：%v。原因是：%v\n", u.Email, err)
		return apperrors.NewInternal()
	}
	return nil
}

func (r *pgUserRepository) FindByID(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	user := &model.User{}

	query := "SELECT * FROM users WHERE uid=$1"

	if err := r.DB.GetContext(ctx, user, query, uid); err != nil {
		return user, apperrors.NewNotFound("uid", uid.String())
	}

	return user, nil
}

func (r *pgUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}

	query := "SELECT * FROM users WHERE email=$1"

	if err := r.DB.GetContext(ctx, user, query, email); err != nil {
		log.Printf("未找到邮箱为 %v 的用户，错误：%v\n", email, err)
		return user, apperrors.NewNotFound("email", email)
	}

	return user, nil
}
