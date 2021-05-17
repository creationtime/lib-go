package fileHandler

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type GetFileContentsRsp struct {
	Content []byte
	Size    int64
	Type    string
	Ext     string
}

//获取根目录下直属所有文件（不包括文件夹及其中的文件）
func GetFile(pathname string) ([]string, error) {
	var res []string
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return res, err
	}

	for _, fi := range rd {
		if !fi.IsDir() {
			fullName := pathname + "/" + fi.Name()
			res = append(res, fullName)
		}
	}
	return res, nil
}

//获取根目录下所有文件（包含文件夹下的文件）
func GetAllFiles(folder string) []string {
	files, _ := ioutil.ReadDir(folder)
	var res []string
	for _, file := range files {
		if file.IsDir() {
			list := GetAllFiles(folder + "/" + file.Name())
			if len(list) > 0 {
				for _, v := range list {
					res = append(res, v)
				}
			}
		} else {
			res = append(res, folder+"/"+file.Name())
		}
	}
	return res
}

//检查文件存在性并创建
func CheckFileIsExistAndCreate(folderPath string) (*os.FileInfo, error) {
	fileInfo, err := CheckFileIsExist(folderPath)
	if err != nil {
		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			return nil, err
		}
	}
	return fileInfo, nil
}

//检查文件存在性
func CheckFileIsExist(folderPath string) (*os.FileInfo, error) {
	fileInfo, err := os.Stat(folderPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("File does not exist.")
		}
	}
	return &fileInfo, nil
}

func CreateFile(newPath string, stream []byte) error {
	f, err := os.Create(newPath)
	if err != nil {
		return err
	}
	if _, err := f.Write(stream); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func GetFileContents(path string) (*GetFileContentsRsp, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	size := int64(binary.Size(content))
	ext := GetFileExt(path)
	_type, err := GetFileType(path, content)
	if err != nil {
		return nil, err
	}
	return &GetFileContentsRsp{
		Content: content,
		Size:    size,
		Type:    _type,
		Ext:     ext,
	}, nil
}

func QuickWriteFile(filePath string, content []byte) error {
	err := ioutil.WriteFile(filePath, content, 0666)
	if err != nil {
		return err
	}
	return nil
}

//skip 循环后退。0 表示调用runtime.Caller()所在的位置，1表示runtime.Caller()所在函数的调用位置，依此类推
func CurrentFile(skip int) (string, error) {
	_, file, _, ok := runtime.Caller(skip)
	if !ok {
		return "", nil
	}
	return file, nil
}

//获取当前执行文件路径
func CurrentExecFileDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", nil
	}
	return dir, nil
}

//获取当前文件路径
func CurrentFileDir(path string) string {
	dir := filepath.Dir(path)
	return dir
}

//打印内容到文件中
func Tracefile(filePath, strContent string) {
	fd, _ := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	fd_time := time.Now().Format("2006-01-02 15:04:05")
	fd_content := strings.Join([]string{"===", fd_time, "===", strContent, "\n"}, "")
	buf := []byte(fd_content)
	fd.Write(buf)
	fd.Close()
}
