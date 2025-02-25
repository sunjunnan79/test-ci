package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUserRoute(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 设置环境变量为 test，加载测试配置文件
	os.Setenv("ENV", "test")
	defer os.Unsetenv("ENV") // 测试结束后清理环境变量

	// 获取初始化的 Gin 路由
	r := setupRouter(InitConfig())

	// 定义测试用例
	tests := []struct {
		name          string
		userID        string
		expectedCode  int
		expectedError string
	}{
		{
			name:          "Valid user",
			userID:        "1", // 假设这个用户存在
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
		{
			name:          "Invalid user ID",
			userID:        "invalid", // 无效的 user_id
			expectedCode:  http.StatusBadRequest,
			expectedError: "Invalid user_id",
		},
		{
			name:          "User not found",
			userID:        "99999", // 假设这个用户不存在
			expectedCode:  http.StatusNotFound,
			expectedError: "User not found",
		},
	}

	// 运行每个测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建请求
			req, _ := http.NewRequest("GET", "/user?user_id="+tt.userID, nil)

			// 创建响应记录器
			w := httptest.NewRecorder()

			// 执行请求
			r.ServeHTTP(w, req)

			// 断言状态码
			assert.Equal(t, tt.expectedCode, w.Code)

			// 如果返回状态不是 200，检查返回的错误信息
			if tt.expectedCode != http.StatusOK {
				var result map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &result)
				assert.NoError(t, err)
				assert.Contains(t, result["error"], tt.expectedError)
			}
		})
	}
}
