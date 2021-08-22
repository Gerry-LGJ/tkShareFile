package Command

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
	"time"
)

/**
 * @brief 发现客户端请求数据错误主动返回错误，调用此函数之后一定要加 return 结束函数，否则可能产生不可预料的后果
 *
 * @param c 请求对象
 * @param httpCode 返回的状态码
 * @param errorInfo 要返回的自定义信息
 */
func RespFail(c *gin.Context, httpCode int, errorInfo string) {
	pc, file, line, _ := runtime.Caller(1) // 返回调用者层的栈信息
	pcName := runtime.FuncForPC(pc).Name()
	errorStack := fmt.Sprintf("%s:%d:%s", file, line, pcName)

	nowTime := time.Now().Local()
	errorTime := fmt.Sprintf("%d/%d/%d - %d:%d:%d", nowTime.Year(), nowTime.Month(), nowTime.Day(),
		nowTime.Hour(), nowTime.Minute(), nowTime.Second())
	debugErrorMsg := "[tkShareFile-error] " + errorTime + " " + errorStack
	c.HTML(httpCode, "error.html",
		gin.H{"title": "Request error", "errorMsg": errorInfo, "debugErrorMsg": debugErrorMsg})
}

/**
 * @brief 捕获致命异常的中间件
 *
 * @return gin.HandlerFunc 中间件处理函数
 */
func FatalError() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.HTML(http.StatusInternalServerError, "error.html",
					gin.H{"title": "Fatal error", "errorMsg": err})
			}
		}()
		c.Next()
	}
}
