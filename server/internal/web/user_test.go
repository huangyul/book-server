package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
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
		// {
		// 	name: "service返回报错",
		// 	mock: func(ctrl *gomock.Controller) service.UserService {
		// 		userService := mocksvc.NewMockUserService(ctrl)
		// 		userService.EXPECT().SignUp(gomock.Any(), "111", "222").Return(int64(1), errors.New("service 返回报错"))
		// 		return userService
		// 	},
		// 	request: Req{
		// 		Username:        "111",
		// 		Password:        "222",
		// 		ConfirmPassword: "222",
		// 	},
		// 	wantCode: 200,
		// 	wantRes: Res{
		// 		Code: 500,
		// 		Msg:  "service 返回报错",
		// 	},
		// },
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
