package handler

import (
	"log"

	"github.com/FuZhouJohn/memrizr/account/model"
	"github.com/FuZhouJohn/memrizr/account/model/apperrors"
	"github.com/gin-gonic/gin"
)

type signupReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

func (h *Handler) Signup(c *gin.Context) {
	var req signupReq

	if ok := bindData(c, &req); !ok {
		return
	}

	u := &model.User{
		Email:    req.Email,
		Password: req.Password,
	}

	err := h.UserService.Signup(c, u)
	if err != nil {
		log.Panicf("注册用户失败：%v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}
}
