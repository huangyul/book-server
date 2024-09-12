package intergration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/huangyul/book-server/internal/intergration/startup"
	"github.com/huangyul/book-server/internal/repository/dao"
	"github.com/huangyul/book-server/internal/web"
	"github.com/stretchr/testify/assert"
)

func TestUser_Edit(t *testing.T) {
	type Req struct {
		Username string `json:"username"`
	}

	db := startup.InitDB()
	server := gin.Default()
	uHdl := startup.InitUserHandler()
	server.Use(func(ctx *gin.Context) {
		ctx.Set("userId", int64(1))
	})
	uHdl.RegisterRoutes(server)

	tests := []struct {
		name     string
		before   func(t *testing.T)
		after    func(t *testing.T)
		req      Req
		wantCode int
		wantRes  web.Result[any]
	}{
		{
			name: "修改成功",
			before: func(t *testing.T) {
				err := db.Model(&dao.User{}).Create(&dao.User{
					Username: "old",
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var user dao.User
				err := db.Model(&dao.User{}).Where("id = ?", 1).First(&user).Error
				assert.NoError(t, err)
				assert.Equal(t, "new", user.Username)
			},
			req: Req{
				Username: "new",
			},
			wantCode: 200,
		},
		{
			name: "用户不存在",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
			},
			req: Req{
				Username: "new",
			},
			wantCode: 200,
			wantRes: web.Result[any]{
				Code: 401002,
				Msg:  "用户不存在",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer db.Exec("truncate table users")

			tc.before(t)
			defer tc.after(t)

			data, err := json.Marshal(tc.req)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/user/edit", bytes.NewReader(data))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantCode, recorder.Code)
			var res web.Result[any]
			assert.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &res))
			assert.Equal(t, tc.wantRes, res)

		})
	}
}
