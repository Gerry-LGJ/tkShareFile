package Preview

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tkShareFile/Command"
)

func Document(conf *Command.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		if !conf.TransportPower.Document || !conf.TransportPower.Download { //检查权限
			Command.RespFail(c, http.StatusBadRequest, "您没有预览文档权限或未开启下载权限，请联系管理员。")
			return
		}
		src, hasKey := c.GetQuery("src")
		if !hasKey {  	//如果没有src参数则返回错误
			Command.RespFail(c, http.StatusBadRequest, "Failed to query 'src' parameter.")
			return
		}

		documentType, err := ExtractFromEnd(src, '.') //获取扩展名
		if err != nil {
			Command.RespFail(c, http.StatusBadRequest, "Failed to resolve file extension.")
			return
		}
		documentType = "." + documentType

		if !SearchStringInArray(conf.DocumentType, documentType) {
			//文件扩展名不在 conf.TextType 则返回错误
			Command.RespFail(c, http.StatusBadRequest, "DocumentType error")
			return
		}

		// TODO 网页预览word等文档，可考虑使用微软的网页预览接口，但需要访问公共资源，故没有好的解决方案
		if documentType == ".pdf" {
			// 参考 https://www.mekau.com/5067.html
			// 参考 https://www.cnblogs.com/kagome2014/p/kagome2014001.html
			c.Redirect(http.StatusMovedPermanently, "/www/pdfJs/web/viewer.html?file=" + src) // 重定向到pdfJs插件中
		} else {
			Command.RespFail(c, http.StatusBadRequest, "File type is not supported.")
			return
		}
	}
}
