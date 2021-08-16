package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	Port int `json:"port"`
	IsDefaultDir bool `json:"is_default_dir"`
	ShareDir string `json:"share_dir"`
	Upload bool `json:"upload"`
	Download bool `json:"download"`
	Player bool `json:"player"`
}

func main() {
	dir, _ := os.Getwd()
	fmt.Println("Current work path: ", dir)
	conf := ReadFileInit()
	r := gin.Default()
	ServerInit(r, conf) 										//服务初始化
	RouteInit(r, conf) 											// 路由初始化

	_ = r.Run(":" + fmt.Sprintf("%v", conf.Port))				// Start Server
}

// 静态资源初始化及添加异常中间件
func ServerInit(r *gin.Engine, conf *Config) {
	r.LoadHTMLGlob("./www/html/*") 						// html模板初始化
	r.Static("/www/img", "./www/img") 			// 静态资源初始化
	r.Static("/www/js", "./www/js") 			// 同上
	r.Static("/www/css", "./www/css") 			// 同上
	r.Static("/www/file", "./www/file") 			// 同上
	r.StaticFS("/shares", http.Dir(conf.ShareDir)) 	// 注册静态文件路径
	r.Use(RespFail())											// 使用中间件捕获异常
}

// 路由初始化，重定向
func RouteInit(r *gin.Engine, conf *Config) {
	r.GET("/", func(c *gin.Context) {
		c.Request.URL.Path = "/share/"
		r.HandleContext(c)
	})
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Request.URL.Path = "/www/img/favicon.ico"
		r.HandleContext(c)
	})
	r.GET("/share/*url", Index(conf))
	r.POST("/upload", Upload(conf))
	r.GET("/player", Player(conf))
}

// ReadFileInit, 读取配置文件// 参数：无// 返回值：//Config结构体信息
func ReadFileInit() *Config {
	data, err := ioutil.ReadFile("./www/config.json")
	CheckError(err, "The search for \"config.json\" failed.")
	conf := Config{}
	err = json.Unmarshal(data, &conf)
	CheckError(err, "\"Config.json\" parsing failed.")

	if !conf.IsDefaultDir {
		fmt.Println("Default folder path is not applicable.")
		fmt.Println("Please enter folder path. Press enter to end the input.")
		for n, err := fmt.Scanln(&conf.ShareDir); n == 0 || err != nil; {
			fmt.Println("Input error, please re-enter. Press enter to end the input.")
			n, err = fmt.Scanln(&conf.ShareDir)
		}
	} else { fmt.Println("Use the default path.") }
	fmt.Println("Share directory : ", conf.ShareDir)
	fmt.Println("Tips: there are still some unknown mistakes (Such as abnormal filenames)" +
		" from development, and if you have doubts, please go to 'https://gitee.com/Time--Chicken' give feedback.")
	return &conf
}

// 返回首页，下载服务
func Index(conf *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Param("url")
		var li string
		list, err := ioutil.ReadDir(conf.ShareDir + url)									// 扫描目录下的所有文件及文件夹
		if err != nil {																		// 因为url访问的可能是一个文件，对错误进行判断
			finfo, err1 := os.Stat(conf.ShareDir + url)										// 读取文件或文件夹的信息
			if err1 != nil {																// 如果都失败了，那路径可定不正确
				//RespFail(c, 500, "打开文件夹失败，请检查文件夹路径")
				CheckError(err1, err1.Error())
				return
			}
			if !finfo.IsDir() && conf.Download {											// 如果打开的是一个文件，并且开放下载权限
				c.Header("Content-Type", "application/octet-stream")
				c.Header("Content-Disposition", "attachment; filename=\"" + finfo.Name() + "\"")
				c.Header("Content-Transfer-Encoding", "binary")
				c.File(conf.ShareDir + url)													// 返回下载文件
			} else { panic("您没有下载权限，请联系管理员。") }								// 不允许下载
			return
		}
		for _, v := range list {
			name := v.Name()
			//li += "<li><a href=\"/share" + url + "/" + name + "\" style=\"color: #fff;\">" + name + "</a></li>\n"
			switch CheckFileType(name) {
			case ".mp3":
				li += strings.Replace("<li><a onclick=\"aTagOnClick('/share" + url + "/" + name + "')\" " +
					"onmouseover=\"aTagOnMouseOver('MP3文件 : 单击下载，双击预览')\" " +
					"onmouseout=\"aTagOnMouseOut()\">" + name +
					"</a></li>\n", "//", "/", 1)
			case ".mp4":
				li += strings.Replace("<li><a onclick=\"aTagOnClick('/share" + url + "/" + name + "')\" " +
					"onmouseover=\"aTagOnMouseOver('MP4文件 : 单击下载，双击预览')\" " +
					"onmouseout=\"aTagOnMouseOut()\">" + name +
					"</a></li>\n", "//", "/", 1)
			default:
				if v.IsDir() {//判断一下是不是文件夹
					li += strings.Replace("<li><a href=\"/share"+url+"/"+name+"\" "+
						"onmouseover=\"aTagOnMouseOver('文件夹 : 单击进入浏览')\" "+
						"onmouseout=\"aTagOnMouseOut()\">"+name+
						"</a></li>\n", "//", "/", 1)
				} else {
					li += strings.Replace("<li><a href=\"/share"+url+"/"+name+"\" "+
						"onmouseover=\"aTagOnMouseOver('普通文件 : 不支持预览，请下载后查看，单击下载')\" "+
						"onmouseout=\"aTagOnMouseOut()\">"+name+
						"</a></li>\n", "//", "/", 1)
				}
			}
		}
		// 封装html
		var filedata []byte
		filedata, err = ioutil.ReadFile("./www/html/index.html")
		CheckError(err, "\"/www/html/index.html\" parsing failed.")
		filestr := string(filedata)
		filestr = strings.Replace(filestr, "{{.path}}", conf.ShareDir + url, -1)
		filestr = strings.Replace(filestr, "{{.li}}", li, -1)
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, filestr)
	}
}

func Player(conf *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !conf.Player || !conf.Download { panic("您没有试听权限或未开启下载权限，请联系管理员。")}// 检查是否开放试听
		src, hasKey := c.GetQuery("src")													// 索取播放路径
		if !hasKey { panic("\"src\" has error") }
		var mediaType string
		switch CheckFileType(src) { //强解析文件类型，有一些奇怪的文件名无法解析
			case ".mp3": mediaType = "audio/mpeg"
			case ".mp4": mediaType = "video/mp4"
			default: panic("file-type error")
		}
		// 开始处理src
		//src = strings.Replace(src, "/share/", "/shares/", 1)
		filename, err := ExtractFromEnd(src, '/')
		if err != nil { CheckError(err, err.Error()) }
		c.HTML(http.StatusOK, "player.html", gin.H{"title": "Player - " + filename, "src": src, "type": mediaType})
	}
}

func CheckError(err error, v interface{}) {
	if err != nil {
		panic(v)
	}
}

func RespFail() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.HTML(400, "error.html", gin.H{"errorMsg": err})
				//c.AbortWithStatusJSON(400, gin.H{"err": err})
			}
		}()
		c.Next()
	}
}

// 可能需要转码，因为有一些文件名很奇怪
func CheckFileType(filename string) string {
	sLen := len(filename)
	for i, _ := range filename { if filename[sLen - i - 1] == '.' { return filename[sLen - i - 1:] } }
	return filename
}

// 可能需要转码，因为有一些文件名很奇怪
func ExtractFromEnd(str string, c byte) (string, error) {
	cLastIndex := strings.LastIndexByte(str, c)
	if cLastIndex >= 0 {
		return str[cLastIndex + 1:], nil
	} else {
		return "", errors.New(fmt.Sprintf("search for '%v' failed", c))
	}
}
 