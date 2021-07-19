package handler

import (
	"log"
	"net/http"

	"github.com/FuZhouJohn/memrizr/account/model"
	"github.com/FuZhouJohn/memrizr/account/model/apperrors"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Me(c *gin.Context) {
	user, exists := c.Get("user")

	if !exists {
		log.Printf("由于未知原因，无法从请求环境中提取用户：%v\n", c)
		err := apperrors.NewInternal()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})

		return
	}

	uid := user.(*model.User).UID

	u, err := h.UserService.Get(c, uid)

	if err != nil {
		log.Printf("无法找到用户:%v\n%v", uid, err)
		e := apperrors.NewNotFound("user", uid.String())

		c.JSON(e.Status(), gin.H{
			"error": e,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": u,
	})
}
