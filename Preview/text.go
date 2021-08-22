package Preview

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
	"tkShareFile/Command"
)

func Text(conf *Command.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !conf.TransportPower.Text || !conf.TransportPower.Download {  //检查权限
			Command.RespFail(c, http.StatusBadRequest, "您没有预览文本权限或未开启下载权限，请联系管理员。")
			return
		}
		src, hasKey := c.GetQuery("src")
		if !hasKey {  													//如果没有src参数则返回错误
			Command.RespFail(c, http.StatusBadRequest, "Failed to query 'src' parameter.")
			return
		}

		textType, err := ExtractFromEnd(src, '.') 					//获取扩展名
		if err != nil {
			Command.RespFail(c, http.StatusBadRequest, "Failed to resolve file extension.")
			return
		}
		textType = "." + textType


		if !SearchStringInArray(conf.TextType, textType) {
			//文件扩展名不在 conf.TextType 则返回错误
			Command.RespFail(c, http.StatusBadRequest, "TextType error")
			return
		}
		var filename string
		filename, err = ExtractFromEnd(src, '/') 					//获取文件名
		if err != nil {
			Command.RespFail(c, http.StatusBadRequest, "Failed to resolve filename.")
			return
		}

		//处理src，解析要打开文本文件的路径
		filePathIndex := strings.Index(src, "://"+c.Request.Host) 		//定位到host位置
		filePathIndex += len("://" + c.Request.Host + "/share")
		filePath := src[filePathIndex:] 								//解析出要访问的文件的路径 /...

		//开始生成网页
		var filedata []byte
		filedata, err = ioutil.ReadFile(conf.ShareDir + filePath)
		if err != nil {
			Command.RespFail(c, http.StatusBadRequest, "\""+conf.ShareDir+filePath+"\" parsing failed.")
			return
		}
		content := string(filedata) //content
		c.HTML(http.StatusOK, "text.html",
			gin.H{"title": "Text - " + filename, "content": content})
	}
}
