package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huangyul/book-server/internal/domain"
	"github.com/huangyul/book-server/internal/pkg/errno"
	"github.com/huangyul/book-server/internal/pkg/jwt"
	mockjwt "github.com/huangyul/book-server/internal/pkg/jwt/mock"
	"github.com/huangyul/book-server/internal/service"
	mocksvc "github.com/huangyul/book-server/internal/service/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserHandler_SignUp(t *testing.T) {
	type Req struct {
		Username        string `json:"username"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	type Res struct {
		Code int
		Msg  string
		Data any
	}
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		request  Req
		wantCode int
		wantRes  Res
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userService := mocksvc.NewMockUserService(ctrl)
				userService.EXPECT().SignUp(gomock.Any(), "111", "222").Return(int64(1), nil)
				return userService
			},
			request: Req{
				Username:        "111",
				Password:        "222",
				ConfirmPassword: "222",
			},
			wantCode: 200,
			wantRes: Res{
				Code: 0,
				Msg:  "",
				Data: float64(1),
			},
		},
		{
			name: "两次密码不一致",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userService := mocksvc.NewMockUserService(ctrl)
				return userService
			},
			request: Req{
				Username:        "111",
				Password:        "222",
				ConfirmPassword: "232",
			},
			wantCode: 200,
			wantRes: Res{
				Code: 400,
				Msg:  "两次密码不一致",
			},
		},
		{
			name: "service返回报错",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userService := mocksvc.NewMockUserService(ctrl)
				userService.EXPECT().SignUp(gomock.Any(), "111", "222").Return(int64(1), errors.New("service 返回报错"))
				return userService
			},
			request: Req{
				Username:        "111",
				Password:        "222",
				ConfirmPassword: "222",
			},
			wantCode: 200,
			wantRes: Res{
				Code: 500,
				Msg:  "service 返回报错",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := tc.mock(ctrl)
			h := NewUserHandler(svc, nil)
			c := gin.Default()
			h.RegisterRoutes(c)
			reqData, err := json.Marshal(tc.request)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/user/signup", bytes.NewBuffer(reqData))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			c.ServeHTTP(recorder, req)
			var res Res
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &res))
			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestUserHandler_Profile(t *testing.T) {
	now := time.Now()
	type UserResp struct {
		ID        float64 `json:"id"`
		Username  string  `json:"username"`
		CreatedAt string  `json:"created_at"`
		UpdatedAt string  `json:"updated_at"`
	}
	type Res struct {
		Code int
		Msg  string
		Data UserResp
	}
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		id       int64
		wantCode int
		wantRes  Res
	}{
		{
			name: "获取成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userService := mocksvc.NewMockUserService(ctrl)
				userService.EXPECT().Profile(gomock.Any(), int64(1)).Return(domain.User{
					ID:        int64(1),
					Username:  "test",
					CreatedAt: now,
					UpdatedAt: now,
				}, nil)
				return userService
			},
			id:       int64(1),
			wantCode: 200,
			wantRes: Res{
				Code: 0,
				Msg:  "",
				Data: UserResp{
					ID:        float64(1),
					Username:  "test",
					CreatedAt: now.Format(time.DateOnly),
					UpdatedAt: now.Format(time.DateOnly),
				},
			},
		},
		{
			name: "service 报错",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userService := mocksvc.NewMockUserService(ctrl)
				userService.EXPECT().Profile(gomock.Any(), int64(1)).Return(domain.User{}, errors.New("service 报错"))
				return userService
			},
			id:       int64(1),
			wantCode: 200,
			wantRes: Res{
				Code: 500,
				Msg:  "service 报错",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := tc.mock(ctrl)
			h := NewUserHandler(svc, nil)
			c := gin.Default()
			h.RegisterRoutes(c)
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/user/profile/%d", tc.id), nil)
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			c.ServeHTTP(recorder, req)
			var res Res
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &res))
			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	type Req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type Token struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	type Res struct {
		Code int
		Msg  string
		Data Token
	}
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (service.UserService, jwt.JWT)
		request  Req
		wantCode int
		wantRes  Res
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, jwt.JWT) {
				userService := mocksvc.NewMockUserService(ctrl)
				userService.EXPECT().Login(gomock.Any(), "111", "222").Return(int64(1), nil)
				jwtMock := mockjwt.NewMockJWT(ctrl)
				jwtMock.EXPECT().GenerateToken(int64(1)).Return("111", "222", nil)
				return userService, jwtMock
			},
			request: Req{
				Username: "111",
				Password: "222",
			},
			wantCode: 200,
			wantRes: Res{
				Code: 0,
				Msg:  "",
				Data: Token{
					AccessToken:  "111",
					RefreshToken: "222",
				},
			},
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) (service.UserService, jwt.JWT) {
				userService := mocksvc.NewMockUserService(ctrl)
				userService.EXPECT().Login(gomock.Any(), "111", "222").Return(int64(0), errno.UserPasswordIncorrect)
				jwtMock := mockjwt.NewMockJWT(ctrl)
				return userService, jwtMock
			},
			request: Req{
				Username: "111",
				Password: "222",
			},
			wantCode: 200,
			wantRes: Res{
				Code: errno.UserNotFound.Code,
				Msg:  "账号或密码不正确",
				Data: Token{},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc, j := tc.mock(ctrl)
			h := NewUserHandler(svc, j)
			c := gin.Default()
			h.RegisterRoutes(c)
			reqData, err := json.Marshal(tc.request)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/user/login", bytes.NewBuffer(reqData))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			c.ServeHTTP(recorder, req)
			assert.Equal(t, tc.wantCode, recorder.Code)
			var res Res
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &res))
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
