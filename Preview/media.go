package Preview

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tkShareFile/Command"
)

func Media(conf *Command.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		if !conf.TransportPower.Media || !conf.TransportPower.Download { // 检查是否开放预览媒体功能
			Command.RespFail(c, http.StatusBadRequest, "您没有权限或未开启下载权限，请联系管理员。")
			return
		}
		src, hasKey := c.GetQuery("src") // 索取播放路径
		if !hasKey {
			Command.RespFail(c, http.StatusBadRequest, "Failed to query 'src' parameter.")
			return
		}
		mediaType, err := ExtractFromEnd(src, '.')
		if err != nil {
			Command.RespFail(c, http.StatusBadRequest, "Failed to match media file type.")
			return
		}
		mediaType = "." + mediaType

		if SearchStringInArray(conf.MediaType.Audio, mediaType) {
			c.Redirect(http.StatusTemporaryRedirect, c.Request.URL.Path + "/audio?src=" + src) 	// 重定向到音频播放模块
		} else if SearchStringInArray(conf.MediaType.Image, mediaType) {
			c.Redirect(http.StatusTemporaryRedirect, c.Request.URL.Path + "/image?src=" + src) 	// 重定向到图像预览模块
		} else if SearchStringInArray(conf.MediaType.Video, mediaType) {
			c.Redirect(http.StatusTemporaryRedirect, c.Request.URL.Path + "/video?src=" + src) 	// 重定向到视频播放模块
		} else {
			Command.RespFail(c, http.StatusBadRequest, "Preview media type is not supported.")//未知媒体类型出错
			return
		}
	}
}

func Audio(conf *Command.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !conf.TransportPower.Media || !conf.TransportPower.Download { // 检查是否开放试听
			Command.RespFail(c, http.StatusBadRequest, "您没有预览音频的权限或未开启下载权限，请联系管理员。")
			return
		}
		src, hasKey := c.GetQuery("src") // 索取播放路径
		if !hasKey {
			Command.RespFail(c, http.StatusBadRequest, "Failed to query 'src' parameter.")
			return
		}

		audioType, err := ExtractFromEnd(src, '.')
		if err != nil {
			Command.RespFail(c, http.StatusBadRequest, "Failed to resolve file extension.")
			return
		}
		audioType = "." + audioType

		if !SearchStringInArray(conf.MediaType.Audio, audioType) {
			//文件扩展名不在 conf.MediaType.Audio 则返回错误
			Command.RespFail(c, http.StatusBadRequest, "MediaType.Audio error")
			return
		}

		// 开始处理src
		//src = strings.Replace(src, "/share/", "/shares/", 1)
		var filename string
		filename, err = ExtractFromEnd(src, '/')
		if err != nil {
			Command.RespFail(c, http.StatusBadRequest, "Failed to resolve filename.")
			return
		}
		// has_bgvid是否视频背景
		c.HTML(http.StatusOK, "player.html",
			gin.H{"title": "Player - " + filename, "src": src, "type": "audio/mpeg", "has_bgvid": "visible"})
	}
}

func Image(conf *Command.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		if !conf.TransportPower.Media || !conf.TransportPower.Download { //检查权限
			Command.RespFail(c, http.StatusBadRequest, "您没有预览图像权限或未开启下载权限，请联系管理员。")
			return
		}
		src, hasKey := c.GetQuery("src")
		if !hasKey { //如果没有src参数则返回错误
			Command.RespFail(c, http.StatusBadRequest, "Failed to query 'src' parameter.")
			return
		}

		imageType, err := ExtractFromEnd(src, '.') //获取扩展名
		if err != nil {
			Command.RespFail(c, http.StatusBadRequest, "Failed to resolve file extension.")
			return
		}
		imageType = "." + imageType

		if !SearchStringInArray(conf.MediaType.Image, imageType) {
			//文件扩展名不在 conf.MediaType.Image 则返回错误
			Command.RespFail(c, http.StatusBadRequest, "MediaType.Image error")
			return
		}

		var filename string
		filename, err = ExtractFromEnd(src, '/') 	//获取文件名
		if err != nil {
			Command.RespFail(c, http.StatusBadRequest,"Failed to resolve filename.")
			return
		}
		c.HTML(http.StatusOK, "image.html", gin.H{"title": "Image - " + filename, "src": src})
	}
}

func Video(conf *Command.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !conf.TransportPower.Media || !conf.TransportPower.Download { // 检查是否开放试听
			Command.RespFail(c, http.StatusBadRequest, "您没有预览视频权限或未开启下载权限，请联系管理员。")
			return
		}
		src, hasKey := c.GetQuery("src") // 索取播放路径
		if !hasKey {
			Command.RespFail(c, http.StatusBadRequest, "Failed to query 'src' parameter.")
			return
		}

		videoType, err := ExtractFromEnd(src, '.')
		if err != nil {
			Command.RespFail(c, http.StatusBadRequest, "Failed to resolve file extension.")
			return
		}
		videoType = "." + videoType

		if !SearchStringInArray(conf.MediaType.Video, videoType) {
			Command.RespFail(c, http.StatusBadRequest, "MediaType.Video error")
			return
		}

		// 开始处理src
		//src = strings.Replace(src, "/share/", "/shares/", 1)
		var filename string
		filename, err = ExtractFromEnd(src, '/')
		if err != nil {
			Command.RespFail(c, http.StatusBadRequest, "Failed to resolve filename.")
			return
		}
		// has_bgvid是否视频背景
		c.HTML(http.StatusOK, "player.html",
			gin.H{"title": "Player - " + filename, "src": src, "type": "video/mp4", "has_bgvid": "hidden"})
	}
}

