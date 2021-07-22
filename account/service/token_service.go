package service

import (
	"context"
	"crypto/rsa"
	"log"

	"github.com/FuZhouJohn/memrizr/account/model"
	"github.com/FuZhouJohn/memrizr/account/model/apperrors"
)

type tokenService struct {
	TokenRepository       model.TokenRepository
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

type TSConfig struct {
	TokenRepository       model.TokenRepository
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

func NewTokenService(c *TSConfig) model.TokenService {
	return &tokenService{
		TokenRepository:       c.TokenRepository,
		PrivKey:               c.PrivKey,
		PubKey:                c.PubKey,
		RefreshSecret:         c.RefreshSecret,
		IDExpirationSecs:      c.IDExpirationSecs,
		RefreshExpirationSecs: c.RefreshExpirationSecs,
	}
}

func (s *tokenService) NewPairFromUser(ctx context.Context, u *model.User, prevTokenID string) (*model.TokenPair, error) {
	idToken, err := generateIDToken(u, s.PrivKey, s.IDExpirationSecs)

	if err != nil {
		log.Printf("为 uid:%v 生成 idToken 时出错，错误：%v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	refreshToken, err := generateRefreshToken(u.UID, s.RefreshSecret, s.RefreshExpirationSecs)

	if err != nil {
		log.Printf("为 uid:%v 生成 refreshToken 时出错，错误：%v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	if err := s.TokenRepository.SetRefreshToken(ctx, u.UID.String(), refreshToken.ID, refreshToken.ExpiresIn); err != nil {
		log.Printf("存储用户 tokenID 时出错，uid：%v。错误：%v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	if prevTokenID != "" {
		if err := s.TokenRepository.DeleteRefreshToken(ctx, u.UID.String(), prevTokenID); err != nil {
			log.Printf("无法删除前一个 refreshToken，uid：%v，tokenID：%v\n", u.UID.String(), prevTokenID)
		}
	}

	return &model.TokenPair{
		IDToken:      idToken,
		RefreshToken: refreshToken.SS,
	}, nil
}
