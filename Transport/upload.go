package Transport

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"tkShareFile/Command"
)

/// 解析多个文件上传中，每个具体的文件的信息
type FileHeader struct {
	ContentDisposition string
	Name               string
	FileName           string ///< 文件名
	ContentType        string
	ContentLength      int64
}

func Upload(conf *Command.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !conf.TransportPower.Upload {
			c.JSON(400, gin.H{"status": 1, "msg": "您没有上传文件的权限，请联系管理员。"})
			return
		}
		uploadDir := c.GetHeader("upload-dir")//检索上传文件指定目录

		//如果在分享目录的根目录，就不加"/"
		if !(len(uploadDir) - strings.LastIndexByte(uploadDir, '/') == 1) { uploadDir += "/" }
		contentLength := c.Request.ContentLength
		if contentLength <= 0 || contentLength > 1024*1024*1024*2 { // 此处貌似是用来限制上传的大小2G
			fmt.Println("Content-Type error")
			return
		}
		ContentType_, hasKey := c.Request.Header["Content-Type"]
		if !hasKey {
			fmt.Println("Content-Type error")
			return
		}
		if len(ContentType_) != 1 {
			fmt.Println("Content-Type count error")
			return
		}
		ContentType := ContentType_[0]
		const BOUNDARY = "; boundary="
		loc := strings.Index(ContentType, BOUNDARY)
		if loc == -1 {
			fmt.Println("Content-Type error, no boundary")
			return
		}
		boundary := []byte(ContentType[(loc + len(BOUNDARY)):])
		fmt.Printf("boundary=[%s]\n\n", boundary)

		readData := make([]byte, 1024*12)
		readTotal := 0
		filename := ""
		for {
			fileHeader, fileData, err := ParseFromHead(readData, readTotal, append(boundary, []byte("\r\n")...), c.Request.Body)
			if err != nil {

				if err.Error() == "reach to sream EOF" {
					//如果是小文件读第一次就返回error则尝试以小文件方式上传，返回307状态码(不改变POST方法重定向)
					c.Redirect(http.StatusTemporaryRedirect, c.Request.URL.Path + "/uploadSmallFile")
				} else {
					fmt.Println(err.Error())
				}
				return
			}
			fmt.Println("File :", fileHeader.FileName)

			//f, err := os.Create(conf.ShareDir + "/" + fileHeader.FileName)//将文件保存在指定目录
			f, err := os.Create(uploadDir + fileHeader.FileName)//将文件保存在指定目录
			if err != nil {
				fmt.Println("create file fail:", err)
				return
			}
			_, err = f.Write(fileData)
			if err != nil {
				fmt.Println("Write \"" + fileHeader.FileName + "\" error.")
				return
			}
			fileData = nil

			//需要反复搜索boundary
			tempData, reachEnd, err := ReadToBoundary(boundary, c.Request.Body, f)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			err = f.Close()
			if err != nil {
				fmt.Println("Close the \"" + fileHeader.FileName + "\" error.")
				return
			}
			if reachEnd {
				filename = fileHeader.FileName //结束时候赋值文件名
				break
			} else {
				copy(readData[0:], tempData)
				readTotal = len(tempData)
				continue
			}
		}

		c.JSON(http.StatusOK, gin.H{"status": 0, "msg": filename + " 上传成功"})
	}
}

/// 解析表单的头部
/// @param read_data 已经从流中读到的数据
/// @param read_total 已经从流中读到的数据长度
/// @param boundary 表单的分割字符串
/// @param stream 输入流
/// @return FileHeader 文件名等信息头
///			[]byte 已经从流中读到的部分
///			error 是否发生错误
func ParseFromHead(readData []byte, readTotal int, boundary []byte, stream io.ReadCloser) (FileHeader, []byte, error) {
	buf := make([]byte, 1024*4)
	foundBoundary := false
	boundaryLoc := -1
	var fileHeader FileHeader
	for {
		readLen, err := stream.Read(buf) //从流中读1024*4个字节
		if err != nil {
			if err != io.EOF {
				return fileHeader, nil, err
			}
			break // TODO 此处有bug，原因是读到了数据同时产生EOF，对于小文件处理不当，解决方案是尝试以小文件的方式上传
		}
		if readTotal+readLen > cap(readData) { //检查读的数据长度是否能放进read_data
			return fileHeader, nil, fmt.Errorf("not found boundary")
		}
		copy(readData[readTotal:], buf[:readLen]) //将stream读到的数据放到read_data
		readTotal += readLen                      //记录读取数据的位置
		if !foundBoundary {                       //该部分在循环中，只执行一次
			boundaryLoc = bytes.Index(readData[:readTotal], boundary) //每从流中读取一次数据时，就检查是否有读到边界
			if -1 == boundaryLoc {                                    //如果没有读到边界，就继续读取
				continue
			}
			foundBoundary = true //如果有读到边界就将found_boundary置为true
		}
		startLoc := boundaryLoc + len(boundary)
		fileHeadLoc := bytes.Index(readData[startLoc:readTotal], []byte("\r\n\r\n"))
		if -1 == fileHeadLoc {
			continue
		}
		fileHeadLoc += startLoc
		ret := false
		fileHeader, ret = ParseFileHeader(readData[startLoc:fileHeadLoc])
		if !ret {
			return fileHeader, nil, fmt.Errorf("ParseFileHeader fail:%s", string(readData[startLoc:fileHeadLoc]))
		}
		return fileHeader, readData[fileHeadLoc+4 : readTotal], nil
	}
	return fileHeader, nil, fmt.Errorf("reach to sream EOF")
}

/// 从流中一直读到文件的末位
/// @return []byte 没有写到文件且又属于下一个文件的数据
/// @return bool 是否已经读到流的末位了
/// @return error 是否发生错误
func ReadToBoundary(boundary []byte, stream io.ReadCloser, target io.WriteCloser) ([]byte, bool, error) {
	readData := make([]byte, 1024*8)
	readDataLen := 0
	buf := make([]byte, 1024*4)
	bLen := len(boundary)
	reachEnd := false
	for !reachEnd {
		readLen, err := stream.Read(buf)
		if err != nil {
			if err != io.EOF && readLen <= 0 {
				return nil, true, err
			}
			reachEnd = true
		}
		//todo: 下面这一句很蠢，值得优化
		copy(readData[readDataLen:], buf[:readLen]) //追加到另一块buffer，仅仅只是为了搜索方便
		readDataLen += readLen
		if readDataLen < bLen+4 {
			continue
		}
		loc := bytes.Index(readData[:readDataLen], boundary)
		if loc >= 0 {
			//找到了结束位置
			_, _ = target.Write(readData[:loc-4])
			return readData[loc:readDataLen], reachEnd, nil
		}

		_, _ = target.Write(readData[:readDataLen-bLen-4])
		copy(readData[0:], readData[readDataLen-bLen-4:])
		readDataLen = bLen + 4
	}
	_, _ = target.Write(readData[:readDataLen])
	return nil, reachEnd, nil
}

/// 解析描述文件信息的头部
/// @return FileHeader 文件名等信息的结构体
/// @return bool 解析成功还是失败
func ParseFileHeader(h []byte) (FileHeader, bool) {
	arr := bytes.Split(h, []byte("\r\n"))
	var outHeader FileHeader
	outHeader.ContentLength = -1
	const (
		CONTENT_DISPOSITION = "Content-Disposition: "
		NAME                = "name=\""
		FILENAME            = "filename=\""
		CONTENT_TYPE        = "Content-Type: "
		CONTENT_LENGTH      = "Content-Length: "
	)
	for _, item := range arr {
		if bytes.HasPrefix(item, []byte(CONTENT_DISPOSITION)) {
			l := len(CONTENT_DISPOSITION)
			arr1 := bytes.Split(item[l:], []byte("; "))
			outHeader.ContentDisposition = string(arr1[0])
			if bytes.HasPrefix(arr1[1], []byte(NAME)) {
				outHeader.Name = string(arr1[1][len(NAME) : len(arr1[1])-1])
			}
			l = len(arr1[2])
			if bytes.HasPrefix(arr1[2], []byte(FILENAME)) && arr1[2][l-1] == 0x22 {
				outHeader.FileName = string(arr1[2][len(FILENAME) : l-1])
			}
		} else if bytes.HasPrefix(item, []byte(CONTENT_TYPE)) {
			l := len(CONTENT_TYPE)
			outHeader.ContentType = string(item[l:])
		} else if bytes.HasPrefix(item, []byte(CONTENT_LENGTH)) {
			l := len(CONTENT_LENGTH)
			s := string(item[l:])
			contentLength, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				log.Printf("content length error:%s", string(item))
				return outHeader, false
			} else {
				outHeader.ContentLength = contentLength
			}
		} else {
			log.Printf("unknown:%s\n", string(item))
		}
	}
	if len(outHeader.FileName) == 0 {
		return outHeader, false
	}
	return outHeader, true
}

/// 此函数将在上传大文件中不能处理小文件时用的函数
/// @param conf 服务基本信息结构体
/// @return gin.HandlerFunc 要处理的函数
func UploadSmallFile(conf *Command.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !conf.TransportPower.Upload {
			c.JSON(400, gin.H{"status": 1, "msg": "您没有上传文件的权限，请联系管理员。"})
			return
		}
		uploadDir := c.GetHeader("upload-dir") 								//检索上传文件指定目录
		if !(len(uploadDir) - strings.LastIndexByte(uploadDir, '/') == 1) { 		//如果在分享目录的根目录，就不加"/"
			uploadDir += "/"
		}
		// 单文件，此处为 Gin 官方的示例代码，唯一不好的地方就是不能上传大文件
		file, _ := c.FormFile("file")
		_ = c.SaveUploadedFile(file, uploadDir + file.Filename) 					// 上传文件至指定目录
		c.JSON(http.StatusOK, gin.H{"status": 0, "msg": file.Filename + " 上传成功"})
	}
}
