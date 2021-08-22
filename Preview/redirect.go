package Preview

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tkShareFile/Command"
)

/**
 * @brief 预览服务重定向初始化
 *
 * @param r 参考Gin的官方代码
 * @param conf 服务信息配置的结构体
 */
func RedirectInit(r *gin.Engine, conf *Command.Config) {

	_preview := r.Group("preview")
	{
		_preview.GET("", Redirect(conf))							//解析出要访问的文件的类型
		_previewMedia := _preview.Group("media") 				//媒体文件
		{
			_previewMedia.GET("", Media(conf))					//解析出要访问的媒体类型
			_previewMedia.GET("/audio", Audio(conf))				//音频
			_previewMedia.GET("/image", Image(conf))				//图像
			_previewMedia.GET("/video", Video(conf))				//视频
		}

		_preview.GET("/text", Text(conf)) 						//文本文件
		_preview.GET("/document", Document(conf))				//文档文件，仅支持pdf
	}
}

/**
 * @brief 重定向解析服务
 *
 * @param conf 服务信息配置的结构体
 * @return gin.HandlerFunc 处理函数
 */
func Redirect(conf *Command.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 大量预览服务时，可考虑重定向到其他链接，以缓解服务压力
		// 根据参数src后的文件后缀判断预览类型，以重定向预览接口
		src, hasKey := c.GetQuery("src")
		if !hasKey {
			Command.RespFail(c, http.StatusBadRequest, "Failed to query 'src' parameter.")
			return
		}
		extractName, err := ExtractFromEnd(src, '.') 										//获取文件后缀名
		if err != nil {
			Command.RespFail(c, http.StatusBadRequest, "Failed to match file type.")	//预览服务重定向解析失败
			return
		}
		extractName = "." + extractName 														//拼接后缀名'.'

		_mediaType := append(conf.MediaType.Audio, conf.MediaType.Image...)
		_mediaType = append(_mediaType, conf.MediaType.Video...)

		if SearchStringInArray(_mediaType, extractName) { 										//如果是媒体文件
			c.Redirect(http.StatusTemporaryRedirect, c.Request.URL.Path + "/media?src=" + src)
		} else if SearchStringInArray(conf.TextType, extractName) {								//如果是文本文件
			c.Redirect(http.StatusTemporaryRedirect, c.Request.URL.Path + "/text?src=" + src)
		} else if SearchStringInArray(conf.DocumentType, extractName) {							//如果是文档文件
			c.Redirect(http.StatusTemporaryRedirect, c.Request.URL.Path + "/document?src=" + src)
		} else {
			Command.RespFail(c, http.StatusBadRequest, "Unknown preview file type.") 	// 无法识别文件类型报错
			return
		}
	}
}
