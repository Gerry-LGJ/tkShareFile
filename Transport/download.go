package Transport

import (
	"github.com/gin-gonic/gin"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"tkShareFile/Command"
	"tkShareFile/Preview"
)

/**
 * @brief 返回首页，下载服务
 *
 * @param conf 服务信息配置的结构体
 * @return gin.HandlerFunc 处理函数
 */
func Index(conf *Command.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Param("url")
		var li string
		list, err := ioutil.ReadDir(conf.ShareDir + url) 		// 扫描目录下的所有文件及文件夹
		if err != nil {                                  		// 因为url访问的可能是一个文件，对错误进行判断
			finfo, err1 := os.Stat(conf.ShareDir + url) 		// 读取文件或文件夹的信息
			if err1 != nil {                            		// 如果都失败了，那路径肯定不正确
				Command.RespFail(c, http.StatusBadRequest,
					"Failed to open folder, please check the folder path.")
				return
			}
			if !finfo.IsDir() && conf.TransportPower.Download { // 如果打开的是一个文件，并且开放下载权限
				c.Header("Content-Type", "application/octet-stream")
				c.Header("Content-Disposition", "attachment; filename=\""+finfo.Name()+"\"")
				c.Header("Content-Transfer-Encoding", "binary")
				c.File(conf.ShareDir + url) 					// 返回下载文件
				return
			} else { 											// 不允许下载
				Command.RespFail(c, http.StatusBadRequest, "您没有下载权限，请联系管理员。")
				return
			}
		}
		for _, v := range list {
			name := v.Name()
			li += AppendLiTagString(conf,Preview.CheckFileType(name),url,name,v)
		}

		// 封装html
		var filedata []byte
		filedata, err = ioutil.ReadFile("./www/html/index.html")
		if err != nil {
			Command.RespFail(c, http.StatusBadRequest, "\"/www/html/index.html\" parsing failed.")
			return
		}
		filestr := string(filedata)
		filestr = strings.Replace(filestr, "{{.title}}", conf.IndexHtmlTitle, -1)
		filestr = strings.Replace(filestr, "{{.path}}", conf.ShareDir+url, -1)
		filestr = strings.Replace(filestr, "{{.li}}", li, -1)
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, filestr)
	}
}

/**
 * @brief 依照配置文件、文件扩展名、访问的url、文件名、文件信息 拼接li标签
 *
 * @param conf 服务信息配置的结构体
 * @param fileExt 文件的扩展名
 * @param url 被访问的路径
 * @param name 文件名
 * @param fileInfo 文件的信息结构体
 * @return string 拼接好的li标签
 */
func AppendLiTagString(conf *Command.Config, fileExt, url, name string, fileInfo fs.FileInfo) string {

	// 先合并媒体文件类型
	doubleClickArray := append(conf.MediaType.Image, conf.MediaType.Audio...)
	doubleClickArray = append(doubleClickArray, conf.MediaType.Video...)
	// 再合并文本文件类型
	doubleClickArray = append(doubleClickArray, conf.TextType...)
	// 后合并文档文件类型
	doubleClickArray = append(doubleClickArray, conf.DocumentType...)
	if Preview.SearchStringInArray(doubleClickArray, fileExt) {          	//判断是否为可预览文件

		return strings.Replace("<li><a onclick=\"aTagOnClick('/share"+url+"/"+name+"')\" "+
			"onmouseover=\"aTagOnMouseOver('Tips : 单击下载，双击预览')\" "+
			"onmouseout=\"aTagOnMouseOut()\">"+name+
			"</a></li>\n", "//", "/", 1)
	} else { 																//如果都不是，则可能为普通文件或目录

		if fileInfo.IsDir() { 												//判断一下是不是文件夹
			return strings.Replace("<li><a href=\"/share"+url+"/"+name+"\" "+
				"onmouseover=\"aTagOnMouseOver('文件夹 : 单击进入浏览')\" "+
				"onmouseout=\"aTagOnMouseOut()\">"+name+
				"</a></li>\n", "//", "/", 1)
		} else {
			return strings.Replace("<li><a href=\"/share"+url+"/"+name+"\" "+
				"onmouseover=\"aTagOnMouseOver('普通文件 : 不支持预览，请下载后查看，单击下载')\" "+
				"onmouseout=\"aTagOnMouseOut()\">"+name+
				"</a></li>\n", "//", "/", 1)
		}
	}
}
