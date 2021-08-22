package Command

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

// 传输权限控制信息结构体
type transportPower struct {
	Upload   bool `json:"upload"`   //允许上传
	Download bool `json:"download"` //允许下载
	Media    bool `json:"media"`    //允许预览 媒体 文件
	Text     bool `json:"text"`     //允许预览 文本 文件
	Document bool `json:"document"` //允许预览 文档 文件
}

// 媒体文件类型分类结构体
type mediaType struct {
	Audio []string `json:"audio"`		//音频
	Image []string `json:"image"`		//图像
	Video []string `json:"video"`		//视频
}

// 服务信息结构体
type Config struct {
	//基本配置
	Port         	int    `json:"port"`           		//端口配置
	IndexHtmlTitle 	string `json:"index_html_title"`	//首页index.html的title标签内容
	IsDefaultDir 	bool   `json:"is_default_dir"` 		//是否指定默认文件夹
	ShareDir     	string `json:"share_dir"`      		//默认文件夹路径，仅在指定默认文件夹时生效，否则需要从键盘输入

	// 权限控制
	TransportPower transportPower `json:"transport_power"`

	//预览文件的文件类型控制
	MediaType      mediaType 	`json:"media_type"`       //媒体文件类型
	TextType       []string 	`json:"text_type"`        //文本文件类型
	DocumentType   []string 	`json:"document_type"`    //文档文件类型
}

/**
 * @brief 读取配置文件
 *
 * @return *Config Config 服务信息配置的结构体
 */
func ReadFileInit() *Config {
	data, err := ioutil.ReadFile("./www/config.json")
	if err != nil {
		panic(err) // The search for \"./www/config.json\" failed.
	}
	conf := Config{}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		panic(err) // "\"Config.json\" parsing failed."
	}

	if !conf.IsDefaultDir {
		fmt.Println("Default folder path is not applicable.")
		fmt.Println("Please enter folder path. Press enter to end the input.")
		for n, err := fmt.Scanln(&conf.ShareDir); n == 0 || err != nil; {
			fmt.Println("Input error, please re-enter. Press enter to end the input.")
			n, err = fmt.Scanln(&conf.ShareDir)
		}
	} else {
		fmt.Println("Use the default path.")
	}
	//fmt.Printf("Share directory : %s\n", conf.ShareDir)
	fmt.Printf("%+v\n", conf)
	fmt.Println("Tips: there are still some unknown mistakes (Such as abnormal filenames)" +
		" from development, and if you have doubts, please go to 'https://gitee.com/Time--Chicken" +
		"/tkShareFile' give feedback.")
	return &conf
}

/**
 * @brief 静态资源初始化及添加异常中间件
 *
 * @param r Default gin.Engine
 * @param conf 服务信息配置的结构体
 */
func ServerInit(r *gin.Engine, conf *Config) {
	r.LoadHTMLGlob("./www/html/*")                 			// html模板初始化
	r.Static("/www/img", "./www/img")              // 静态资源初始化
	r.Static("/www/js", "./www/js")                // 同上
	r.Static("/www/css", "./www/css")              // 同上
	r.Static("/www/file", "./www/file")            // 同上
	r.StaticFS("/www/pdfJs", http.Dir("./www/pdfjs-2.9.359-dist")) // 注册pdfJs插件资源路径
	r.StaticFS("/shares", http.Dir(conf.ShareDir)) 		// 注册静态文件路径
	r.Use(FatalError())                             		 		// 使用中间件捕获致命异常
}
