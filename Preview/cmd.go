package Preview

import (
	"errors"
	"fmt"
	"strings"
)

/**
 * @brief 检查文件类型，仅依靠后缀名判断，与ExtractFromEnd(filename, '.') 相似，多出了错误判断
 *
 * @param filename 文件名
 * @return string 文件扩展名，例如:".mp3"
 */
func CheckFileType(filename string) string {
	sLen := len(filename)
	for i, _ := range filename {
		if filename[sLen-i-1] == '.' {
			return filename[sLen-i-1:]
		}
	}
	return filename
}

/**
 * @brief 获取c指定字符后的字符串，不包括c自己
 *
 * @param str 要查找的字符串
 * @param c 要查找的字符
 * @return string 找到的字符串，否则error不为空
 * @return error 错误信息
 */
func ExtractFromEnd(str string, c byte) (string, error) {
	cLastIndex := strings.LastIndexByte(str, c)
	if cLastIndex >= 0 {
		return str[cLastIndex+1:], nil
	} else {
		return "", errors.New(fmt.Sprintf("search for '%v' failed", c))
	}
}

/**
 * @brief 检查a []string中是否存在x string 是返回true, 否则false
 *
 * @param a 要查找的字符串数组
 * @param x 要查找的字符串
 * @return true 存在
 * @return false 不存在
 */
func SearchStringInArray(a []string, x string) bool {
	/*sort.Strings(a)                    //排序
	StrPos := sort.SearchStrings(a, x) //二分查找，会出现audio的.mp3变成image的.jpg 不知道为啥
	if StrPos != len(a) {
		if x == a[StrPos] {
			return true
		} else {
			return false
		}
	} else {
		return false
	}*/
	for _, value := range a { // 当判断类型过多时就需要二分查找了
		if value == x {
			return true
		} else {
			continue
		}
	}
	return false
}
