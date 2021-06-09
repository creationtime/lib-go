package fileHandler

import (
	"bytes"
	"encoding/hex"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var fileTypeMap sync.Map

func init() {
	fileTypeMap.Store("ffd8ffe0", "jpg")              //JPEG (jpg)
	fileTypeMap.Store("89504e47", "png")              //PNG (png)
	fileTypeMap.Store("47494638", "gif")              //GIF (gif)
	fileTypeMap.Store("49492a00", "tif")              //TIFF (tif)
	fileTypeMap.Store("424d228c010000000000", "bmp")  //16色位图(bmp)
	fileTypeMap.Store("424d8240090000000000", "bmp")  //24位位图(bmp)
	fileTypeMap.Store("424d8e1b030000000000", "bmp")  //256色位图(bmp)
	fileTypeMap.Store("41433130", "dwg")              //CAD (dwg)
	fileTypeMap.Store("3c21444f435459504520", "html") //HTML (html)   3c68746d6c3e0  3c68746d6c3e0
	fileTypeMap.Store("3c68746d6c3e", "html")         //HTML (html)   3c68746d6c3e0  3c68746d6c3e0
	fileTypeMap.Store("3c21646f637479706520", "htm")  //HTM (htm)
	fileTypeMap.Store("48544d4c207b0d0a0942", "css")  //css
	fileTypeMap.Store("696b2e71623d696b2e71", "js")   //js
	fileTypeMap.Store("7b5c727466", "rtf")            //Rich Text Format (rtf)
	fileTypeMap.Store("38425053", "psd")              //Photoshop (psd)
	fileTypeMap.Store("46726f6d3a203d3f6762", "eml")  //Email [Outlook Express 6] (eml)
	fileTypeMap.Store("d0cf11e0a1b11ae10000", "doc")  //MS Excel 注意：word、msi 和 excel的文件头一样
	fileTypeMap.Store("d0cf11e0a1b11ae10000", "vsd")  //Visio 绘图
	fileTypeMap.Store("5374616E64617264204A", "mdb")  //MS Access (mdb)
	fileTypeMap.Store("252150532D41646F6265", "ps")
	fileTypeMap.Store("255044462d312e350d0a", "pdf")  //Adobe Acrobat (pdf)
	fileTypeMap.Store("2e524d46000000120001", "rmvb") //rmvb/rm相同
	fileTypeMap.Store("464c5601050000000900", "flv")  //flv与f4v相同
	fileTypeMap.Store("00000020667479706d70", "mp4")
	fileTypeMap.Store("49443303000000002176", "mp3") //49443303
	fileTypeMap.Store("000001ba210001000180", "mpg") //
	fileTypeMap.Store("3026b2758e66cf11a6d9", "wmv") //wmv与asf相同
	fileTypeMap.Store("52494646e27807005741", "wav") //Wave (wav)
	fileTypeMap.Store("52494646d07d60074156", "avi")
	fileTypeMap.Store("4d546864000000060001", "mid") //MIDI (mid)
	fileTypeMap.Store("504b0304140000000800", "zip")
	fileTypeMap.Store("526172211a0700cf9073", "rar")
	fileTypeMap.Store("235468697320636f6e66", "ini")
	fileTypeMap.Store("504b03040a0000000000", "jar")
	fileTypeMap.Store("4d5a9000030000000400", "exe")        //可执行文件
	fileTypeMap.Store("3c25402070616765206c", "jsp")        //jsp文件
	fileTypeMap.Store("4d616e69666573742d56", "mf")         //MF文件
	fileTypeMap.Store("3c3f786d6c", "xml")                  //xml文件
	fileTypeMap.Store("494e5345525420494e54", "sql")        //xml文件
	fileTypeMap.Store("7061636b616765207765", "java")       //java文件
	fileTypeMap.Store("406563686f206f66660d", "bat")        //bat文件
	fileTypeMap.Store("1f8b0800000000000000", "gz")         //gz文件
	fileTypeMap.Store("6c6f67346a2e726f6f74", "properties") //bat文件
	fileTypeMap.Store("cafebabe0000002e0041", "class")      //bat文件
	fileTypeMap.Store("49545346030000006000", "chm")        //bat文件
	fileTypeMap.Store("04000000010000001300", "mxp")        //bat文件
	fileTypeMap.Store("504b0304140006000800", "docx")       //docx文件
	fileTypeMap.Store("d0cf11e0a1b11ae10000", "wps")        //WPS文字wps、表格et、演示dps都是一样的
	fileTypeMap.Store("6431303a637265617465", "torrent")
	fileTypeMap.Store("6D6F6F76", "mov")                       //Quicktime (mov)
	fileTypeMap.Store("FF575043", "wpd")                       //WordPerfect (wpd)
	fileTypeMap.Store("CFAD12FEC5FD746F", "dbx")               //Outlook Express (dbx)
	fileTypeMap.Store("2142444E", "pst")                       //Outlook (pst)
	fileTypeMap.Store("AC9EBD8F", "qdf")                       //Quicken (qdf)
	fileTypeMap.Store("E3828596", "pwl")                       //Windows Password (pwl)
	fileTypeMap.Store("2E7261FD", "ram")                       //Real Audio (ram)
	fileTypeMap.Store("1a45dfa39f4286810142", "webm")          //Real Audio (audio/webm;codecs=opus)
	fileTypeMap.Store("fff14c40039ffcde0200", "aac")           //Real Audio (audio/x-aac)
	fileTypeMap.Store("44656C69766572792D646174653A", "email") //Email [thorough only] (eml)
}

//获取文件类型
func GetFileType(path string, fSrc []byte) (string, error) {
	var fileType string
	fileCode, err := bytesToHexString(fSrc)
	if err != nil {
		return "", err
	}

	file, _ := CurrentFile(1)
	Tracefile(CurrentFileDir(file)+"/type_record", path+" ===> "+fileCode[:50])

	fileTypeMap.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		if strings.HasPrefix(fileCode, strings.ToLower(k)) || strings.HasPrefix(k, strings.ToLower(fileCode)) {
			fileType = v
			return false
		}
		return true
	})
	return fileType, nil
}

// 获取前面结果字节的二进制
func bytesToHexString(src []byte) (string, error) {
	res := bytes.Buffer{}
	if src == nil || len(src) <= 0 {
		return "", nil
	}
	temp := make([]byte, 0)
	for _, v := range src {
		sub := v & 0xFF
		hv := hex.EncodeToString(append(temp, sub))
		if len(hv) < 2 {
			if _, err := res.WriteString(strconv.FormatInt(int64(0), 10)); err != nil {
				return "", err
			}
		}
		if _, err := res.WriteString(hv); err != nil {
			return "", err
		}
	}
	return res.String(), nil
}

/*
做 Web 应用程序时，经常需要对用户上传的文件类型做一下检查，比如判断上传的是否是 png 、gif、jpg 等图片类型，还是 pdf。并针对不同的类型做一些处理，比如在需要图片的场合，如果上传的文件类型为非图片，那么就要拒绝并告诉用户需要一张图片。

这种需求，往往我们都是通过上传的文件的扩展名来判断，比如如果以 .jpg 结尾，那么我们可能就认为是 jpg 图片。

但这种判断方式往往存在一个隐患，就是恶意用户可能会把一些恶意程序简单的以 .jpg 结尾来骗过我们的逻辑。在这种情况下，最好的方式就是根据上传内容的前几个字节来判断。

根据上传的内容来判断上传的内容类型，一般情况前三个字节或者前八个自己就可以了。比如 PNG 格式的图片，以十六进制 89504E47 开头
*/
func GetFileContentType(out *os.File) (string, error) {
	// 只需要前 512 个字节就可以了
	buffer := make([]byte, 512)
	_, err := out.Read(buffer)

	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(buffer)
	return contentType, nil
}

//获取文件后缀格式
func GetFileExt(fileName string) string {
	return strings.ToLower(filepath.Ext(fileName))
}
