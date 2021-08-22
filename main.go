package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"time"
	"tkShareFile/Command"
	"tkShareFile/Preview"
	"tkShareFile/Transport"
)

func main() {
	defer func() {											// 全局捕获异常输出，延时 10s 结束程序。
		if err := recover(); err != nil {
			fmt.Println(err)
			time.Sleep(time.Second*10)
		}
	}()
	dir, _ := os.Getwd()
	fmt.Println("Current work path: ", dir)				// 显示工作目录，以便调试
	conf := Command.ReadFileInit()							// 通过配置文件初始化服务信息

	// Gin 日志初始化
	gin.DisableConsoleColor() 								// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
	f, _ := os.Create("./www/tkShareFile.log") 		// 记录到文件
	gin.DefaultWriter = io.MultiWriter(f)
	fmt.Printf("For more output information, please visit '%s\\www\\tkShareFile.log'.\n", dir)

	r := gin.Default()
	Command.ServerInit(r, conf) 							//服务初始化
	RouteInit(r, conf)  									// 路由初始化
	fmt.Println("After service initialization, starting the service ......")
	fmt.Printf("You can test the service through 'http://localhost:%d'\n", conf.Port)
	fmt.Println("Type 'Ctrl + C' to terminate the service.")
	_ = r.Run(":" + fmt.Sprintf("%v", conf.Port)) 	// Start Server
}

/**
 * @brief 路由初始化，重定向
 *
 * @param r *gin.Engine
 * @param conf 服务基本信息结构体
 */
func RouteInit(r *gin.Engine, conf *Command.Config) {

	r.GET("", func(c *gin.Context) {					// 主页重定向
		c.Redirect(http.StatusFound, "/share/")
	})
	r.GET("/favicon.ico", func(c *gin.Context) { 	// 图标重定向
		c.Request.URL.Path = "/www/img/favicon.ico"
		r.HandleContext(c)
	})
	r.GET("/share/*url", Transport.Index(conf))		// 注册首页服务

	_upload := r.Group("upload")
	{
		_upload.POST("", Transport.Upload(conf))		// 注册上传服务
		_upload.POST("/uploadSmallFile", Transport.UploadSmallFile(conf)) // 注册上传小文件服务
	}

	Preview.RedirectInit(r, conf) 								// 预览服务重定向服务初始化
}

