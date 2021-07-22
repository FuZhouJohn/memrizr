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

func TestSignup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("账号和密码为必填项", func(t *testing.T) {
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.User")).Return(nil)

		rr := httptest.NewRecorder()

		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		reqBody, err := json.Marshal(gin.H{
			"email": "",
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockUserService.AssertNotCalled(t, "Signup")
	})

	t.Run("邮箱格式错误", func(t *testing.T) {
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.User")).Return(nil)

		rr := httptest.NewRecorder()

		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		reqBody, err := json.Marshal(gin.H{
			"email":    "hello@world",
			"password": "testpassword",
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockUserService.AssertNotCalled(t, "Signup")
	})

	t.Run("密码太短", func(t *testing.T) {
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.User")).Return(nil)

		rr := httptest.NewRecorder()

		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		reqBody, err := json.Marshal(gin.H{
			"email":    "hello@world.com",
			"password": "test",
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockUserService.AssertNotCalled(t, "Signup")
	})

	t.Run("密码太长", func(t *testing.T) {
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.User")).Return(nil)

		rr := httptest.NewRecorder()

		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		reqBody, err := json.Marshal(gin.H{
			"email":    "hello@world.com",
			"password": "testpasswordtestpasswordtestpas",
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockUserService.AssertNotCalled(t, "Signup")
	})

	t.Run("调用 UserService 出错", func(t *testing.T) {
		u := &model.User{
			Email:    "hello@world.com",
			Password: "testpassword",
		}

		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*context.emptyCtx"), u).Return(apperrors.NewConflict("用户已经存在", u.Email))

		rr := httptest.NewRecorder()

		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		reqBody, err := json.Marshal(gin.H{
			"email":    u.Email,
			"password": u.Password,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusConflict, rr.Code)
		mockUserService.AssertExpectations(t)
	})

	t.Run("创建令牌成功", func(t *testing.T) {
		u := &model.User{
			Email:    "hello@world.com",
			Password: "testpassword",
		}

		mockTokenResp := &model.TokenPair{
			IDToken:      "IDToken",
			RefreshToken: "RefreshToken",
		}

		mockUserService := new(mocks.MockUserService)
		mockTokenService := new(mocks.MockTokenService)

		mockUserService.On("Signup", mock.AnythingOfType("*context.emptyCtx"), u).Return(nil)
		mockTokenService.On("NewPairFromUser", mock.AnythingOfType("*context.emptyCtx"), u, "").Return(mockTokenResp, nil)

		rr := httptest.NewRecorder()

		router := gin.Default()

		NewHandler(&Config{
			R:            router,
			UserService:  mockUserService,
			TokenService: mockTokenService,
		})

		reqBody, err := json.Marshal(gin.H{
			"email":    u.Email,
			"password": u.Password,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"tokens": mockTokenResp,
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)
	})

	t.Run("创建令牌失败", func(t *testing.T) {
		u := &model.User{
			Email:    "hello@world.com",
			Password: "testpassword",
		}

		mockErrorResponse := apperrors.NewInternal()

		mockUserService := new(mocks.MockUserService)
		mockTokenService := new(mocks.MockTokenService)

		mockUserService.On("Signup", mock.AnythingOfType("*context.emptyCtx"), u).Return(nil)
		mockTokenService.On("NewPairFromUser", mock.AnythingOfType("*context.emptyCtx"), u, "").Return(nil, mockErrorResponse)

		rr := httptest.NewRecorder()

		router := gin.Default()

		NewHandler(&Config{
			R:            router,
			UserService:  mockUserService,
			TokenService: mockTokenService,
		})

		reqBody, err := json.Marshal(gin.H{
			"email":    u.Email,
			"password": u.Password,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockErrorResponse,
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)
	})
}
