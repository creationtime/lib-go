package fileHandler

import (
	"encoding/binary"
	"io/ioutil"
	"os"

	. "github.com/creationtime/lib-go/fileHandler/check"
)

type GetFileContentsRsp struct {
	Content []byte
	Size    int64
	Type    string
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

//检查文件存在性
func CheckFileIsExistAndCreate(folderPath string) (string, error) {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			return "", err
		}
	}
	return folderPath, nil
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
	_type, err := GetFileType(content)
	if err != nil {
		return nil, err
	}
	return &GetFileContentsRsp{
		Content: content,
		Size:    size,
		Type:    _type,
	}, nil
}
