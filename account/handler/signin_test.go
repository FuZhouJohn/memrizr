package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FuZhouJohn/memrizr/account/model"
	"github.com/FuZhouJohn/memrizr/account/model/apperrors"
	"github.com/FuZhouJohn/memrizr/account/model/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserService := new(mocks.MockUserService)
	mockTokenService := new(mocks.MockTokenService)

	router := gin.Default()

	NewHandler(&Config{
		R:            router,
		UserService:  mockUserService,
		TokenService: mockTokenService,
	})

	t.Run("BadRequest", func(t *testing.T) {
		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"email":    "bademail",
			"password": "short",
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockUserService.AssertNotCalled(t, "Signin")
		mockTokenService.AssertNotCalled(t, "NewTokensFromUser")
	})

	t.Run("UserService.Signin 执行失败", func(t *testing.T) {
		email := "zhuang@test.com"
		password := "testpassword"

		mockUSArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&model.User{
				Email:    email,
				Password: password,
			},
		}
		mockError := apperrors.NewAuthorization("用户名或密码错误")

		mockUserService.On("Signin", mockUSArgs...).Return(mockError)

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"email":    email,
			"password": password,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		mockUserService.AssertCalled(t, "Signin", mockUSArgs...)
		mockTokenService.AssertNotCalled(t, "NewTokensFromUser")
	})
}
