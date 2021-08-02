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

	t.Run("创建 Token 成功", func(t *testing.T) {
		email := "zhuang@test.com"
		password := "testpassword2"

		mockUSArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&model.User{Email: email, Password: password},
		}

		mockUserService.On("Signin", mockUSArgs...).Return(nil)

		mockTSArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&model.User{
				Email:    email,
				Password: password,
			},
			"",
		}
		mockTokenPair := &model.TokenPair{
			IDToken:      "idToken",
			RefreshToken: "refreshToken",
		}

		mockTokenService.On("NewPairFromUser", mockTSArgs...).Return(mockTokenPair, nil)

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

		respBody, err := json.Marshal(gin.H{
			"tokens": mockTokenPair,
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockUserService.AssertCalled(t, "Signin", mockUSArgs...)
		mockTokenService.AssertCalled(t, "NewPairFromUser", mockTSArgs...)
	})

	t.Run("创建 Token 失败", func(t *testing.T) {
		email := "zhuang@test.com"
		password := "testpassword3"

		mockUSArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&model.User{Email: email, Password: password},
		}

		mockUserService.On("Signin", mockUSArgs...).Return(nil)

		mockTSArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&model.User{
				Email:    email,
				Password: password,
			},
			"",
		}
		mockError := apperrors.NewInternal()

		mockTokenService.On("NewPairFromUser", mockTSArgs...).Return(nil, mockError)

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

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockUserService.AssertCalled(t, "Signin", mockUSArgs...)
		mockTokenService.AssertCalled(t, "NewPairFromUser", mockTSArgs...)
	})
}
