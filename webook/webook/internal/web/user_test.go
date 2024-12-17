package web

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"webook/webook/internal/service"
	svcmocks "webook/webook/internal/service/mocks"
)

func TestUserHandler_SignUpV2(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserServicer
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserServicer {
				userSvc := svcmocks.NewMockUserServicer(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).
					Return(nil)
				return userSvc
			},
			reqBody: `
					{
						"email":"1234@qq.com",
						"password":"hello#world123",
						"confirm_password":"hello#world123"
					}
			`,
			wantCode: 200,
			wantBody: "welcome",
		},
		{
			name: "参数不对，bind 失败",
			mock: func(ctrl *gomock.Controller) service.UserServicer {
				userSvc := svcmocks.NewMockUserServicer(ctrl)
				return userSvc
			},

			reqBody: `
			{
				"email": "123@qq.com",
				"password": "hello#world123"
			}
			`,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gin.SetMode(gin.DebugMode)
			server := gin.New()
			server.Use(gin.Logger(), gin.Recovery())
			hdl := NewUserHandler(tc.mock(ctrl), nil)
			hdl.RegisterRoutes(server)

			req, _ := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			req.Header.Set("Content-Type", "application/json")
			recoder := httptest.NewRecorder()
			server.ServeHTTP(recoder, req)
			assert.Equal(t, tc.wantCode, recoder.Code)
			assert.Equal(t, tc.wantBody, recoder.Body.String())
		})
	}
}

func TestUserHandler_SignUp_Simplified(t *testing.T) {
	//reqBody := `{"email":"12345678@qq.com","password":"hello#world123","confirm_password":"hello#world123"}`
	//wantCode := 200
	//wantBody := "welcome"
	reqBody := `
			{
				"email": "123@qq.com",
				"password": "hello#world123"
			}
			`
	wantCode := http.StatusBadRequest
	wantBody := ""

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userSvc := svcmocks.NewMockUserServicer(ctrl)
	userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)

	server := gin.Default()
	hdl := NewUserHandler(userSvc, nil)
	hdl.RegisterRoutes(server)

	req, _ := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, req)

	assert.Equal(t, wantCode, recorder.Code)
	assert.Equal(t, wantBody, recorder.Body.String())
}

func TestUserHandler_SignUpV1(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserServicer
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserServicer {
				userSvc := svcmocks.NewMockUserServicer(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).
					Return(nil)
				return userSvc
			},
			reqBody: `
					{
						"email":"1234@qq.com",
						"password":"hello#world123",
						"confirm_password":"hello#world123"
					}
			`,
			wantCode: 200,
			wantBody: "welcome",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gin.SetMode(gin.DebugMode)
			server := gin.New()
			server.Use(gin.Logger(), gin.Recovery())
			hdl := NewUserHandler(tc.mock(ctrl), nil)
			hdl.RegisterRoutes(server)

			req, _ := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			req.Header.Set("Content-Type", "application/json")
			recoder := httptest.NewRecorder()
			server.ServeHTTP(recoder, req)
			assert.Equal(t, tc.wantCode, recoder.Code)
			assert.Equal(t, tc.wantBody, recoder.Body.String())
		})
	}
}
