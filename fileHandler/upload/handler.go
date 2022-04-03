package upload

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	. "github.com/lights-T/lib-go/fileHandler/upload/config"
	"github.com/lights-T/lib-go/logger"
	"github.com/nfnt/resize"
)

const FIELD = "files"

// 支持的图片后缀名
var supportImageExtNames = []string{".jpg", ".jpeg", ".png", ".ico", ".svg", ".bmp", ".gif"}

func NewUploader(confPath string) (upHandler *UpHandler, err error) {
	conf, err := NewConfig(confPath)
	if err != nil {
		return upHandler, err
	}
	upHandler = new(UpHandler)
	upHandler.Conf = conf
	return upHandler, nil
}

/**
check a file is a image or not
*/
func isImage(extName string) bool {
	for i := 0; i < len(supportImageExtNames); i++ {
		if supportImageExtNames[i] == extName {
			return true
		}
	}
	return false
}

/**
Generate thumbnail
*/
func (h *UpHandler) thumbnail(imagePath string) (outputPath string, err error) {
	var (
		file     *os.File
		img      image.Image
		filename = path.Base(imagePath)
	)

	extname := strings.ToLower(path.Ext(imagePath))

	timeDir, err := mkdirMonthDir(MkdirThumbnail)
	if err != nil {
		return outputPath, errors.New(fmt.Sprintf("mkdirMonthDir a failure, err: %s", err.Error()))
	}
	outputPath = path.Join(h.Conf.Upload.Path, h.Conf.Upload.Image.Path, h.Conf.Upload.Image.Thumbnail.Path, timeDir, filename)

	// Read the file
	if file, err = os.Open(imagePath); err != nil {
		return outputPath, errors.New(fmt.Sprintf("Open a failure, err: %s", err.Error()))
	}

	defer file.Close()

	// decode jpeg into image.Image
	switch extname {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
		break
	case ".png":
		img, err = png.Decode(file)
		break
	case ".gif":
		img, err = gif.Decode(file)
		break
	default:
		return outputPath, errors.New(fmt.Sprintf("UnSupport file type" + extname))
	}

	if img == nil {
		return outputPath, errors.New(fmt.Sprintf("Generate thumbnail fail..."))
	}

	m := resize.Thumbnail(uint(h.Conf.Upload.Image.Thumbnail.MaxWidth), uint(h.Conf.Upload.Image.Thumbnail.MaxHeight), img, resize.Lanczos3)

	out, err := os.Create(outputPath)
	if err != nil {
		return outputPath, errors.New(fmt.Sprintf("Create a failure, err: %s", err.Error()))
	}
	defer out.Close()

	// write new image to file

	//decode jpeg/png/gif into image.Image
	switch extname {
	case ".jpg", ".jpeg":
		jpeg.Encode(out, m, nil)
		break
	case ".png":
		png.Encode(out, m)
		break
	case ".gif":
		gif.Encode(out, m, nil)
		break
	default:
		return outputPath, errors.New(fmt.Sprintf("UnSupport file type" + extname))
	}

	return
}

func mkdirMonthDir(mkType int32) (dirName string, err error) {
	y, m, _ := time.Now().Date()
	timeDir := fmt.Sprintf("%d%d", y, m)
	switch mkType {
	case MkdirImage:
		if err = os.MkdirAll(path.Join(Config.Upload.Path, Config.Upload.Image.Path, timeDir), PermFile); err != nil {
			logger.Errorf("Failed to create the picture directory, err: %s", err.Error())
			return
		}
	case MkdirFile:
		if err = os.MkdirAll(path.Join(Config.Upload.Path, Config.Upload.File.Path, timeDir), PermFile); err != nil {
			logger.Errorf("Failed to create the file directory, err: %s", err.Error())
			return
		}
	case MkdirThumbnail:
		if err = os.MkdirAll(path.Join(Config.Upload.Path, Config.Upload.Image.Path, Config.Upload.Image.Thumbnail.Path, timeDir), PermFile); err != nil {
			logger.Errorf("Failed to create the thumbnail directory, err: %s", err.Error())
			return
		}
	}
	return timeDir, nil
}

type UpHandler struct {
	Conf *ConfigType
}

type imageRspInfo struct {
	Hash     string
	FileName string
	Origin   string
	Size     int64
}

// UploaderImage Upload image handler
func (h *UpHandler) UploaderImage(b []byte) (rsp *imageRspInfo, err error) {
	var (
		maxUploadSize = h.Conf.Upload.Image.MaxSize // 最大上传大小
		distPath      string                        // 最终的输出目录
		src           multipart.File
		dist          *os.File
	)

	rsp = &imageRspInfo{}

	// Source
	file := &multipart.FileHeader{}
	err = json.Unmarshal(b, file)

	extname := strings.ToLower(path.Ext(file.Filename))

	if isImage(extname) == false {
		return rsp, errors.New("UnSupport upload file type " + extname)
	}

	if file.Size > int64(maxUploadSize) {
		return rsp, errors.New("Upload file too large, The max upload limit is " + strconv.Itoa(int(maxUploadSize)))
	}

	if src, err = file.Open(); err != nil {
	}
	defer src.Close()

	hash := md5.New()

	io.Copy(hash, src)

	md5string := hex.EncodeToString(hash.Sum([]byte("")))
	fileName := md5string + extname
	rsp.Hash = md5string
	rsp.FileName = fileName
	rsp.Origin = file.Filename
	rsp.Size = file.Size

	timeDir, err := mkdirMonthDir(MkdirImage)
	if err != nil {
		return rsp, errors.New(fmt.Sprintf("Failed to mkdirMonthDir, err: %s", err.Error()))
	}

	// Destination
	distPath = path.Join(h.Conf.Upload.Path, h.Conf.Upload.Image.Path, timeDir, fileName)
	if dist, err = os.Create(distPath); err != nil {
		return rsp, errors.New(fmt.Sprintf("Create a failure, err: %s", err.Error()))
	}

	defer dist.Close()

	// FIXME: open 2 times
	if src, err = file.Open(); err != nil {
		return rsp, errors.New(fmt.Sprintf("Open a failure, err: %s", err.Error()))
	}

	// Copy
	io.Copy(dist, src)

	// 压缩缩略图
	// 不管成功与否，都会进行下一步的返回
	if _, err = h.thumbnail(distPath); err != nil {
		return rsp, errors.New(fmt.Sprintf("Failed to create thumbnail, err: %s", err.Error()))
	}

	return rsp, nil
}
