package config

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type ConfigType struct {
	Http   HttpConfig
	Upload TConfig
	Env    string
	Mode   string
}

var Config ConfigType

type FileConfig struct {
	Path      string   `valid:"required,length(1|20)"`  // 普通文件的存放目录
	MaxSize   int      `valid:"required"`               // 普通文件上传的限制大小，单位byte, 最大单位1GB
	AllowType []string `valid:"required,length(0|100)"` // 允许上传的文件后缀名
}

type ImageConfig struct {
	Path      string `valid:"required,length(1|20)"` // 图片存储路径
	MaxSize   int    `valid:"required"`              // 最大图片上传限制，单位byte
	Thumbnail ThumbnailConfig
}

type ThumbnailConfig struct {
	Path      string `valid:"required,length(1|20)"` // 缩略图存放路径
	MaxWidth  int    `valid:"required"`              // 缩略图最大宽度
	MaxHeight int    `valid:"required"`              // 缩略图最大高度
}

type TConfig struct {
	Path      string `valid:"required,length(1|20)"` //文件上传的根目录
	UrlPrefix string `valid:"required,length(0|20)"` // api的url前缀
	File      FileConfig
	Image     ImageConfig
}

type Uploader struct {
	Upload   *gin.RouterGroup
	Download *gin.RouterGroup
	Config   TConfig
	Engine   *gin.Engine
}

const (
	//configFilePath = "./files/upload/config/config.yaml"
	PermFile       = 0755
	MkdirImage     = 1
	MkdirFile      = 2
	MkdirThumbnail = 3
)

func NewConfig(configFilePath string) (config *ConfigType, err error) {
	var (
		configFile string
		yamlFile   []byte
	)

	configFile, err = filepath.Abs(configFilePath)
	if err != nil {
		return &Config, err
	}

	yamlFile, err = ioutil.ReadFile(configFile)
	if err != nil {
		return &Config, err
	}

	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		return &Config, err
	}

	Config.Env = os.Getenv("GO_ENV")
	if Config.Env == gin.ReleaseMode || Config.Env == "production" || Config.Env == "publish" {
		Config.Mode = gin.ReleaseMode
	} else if Config.Env == gin.TestMode {
		Config.Mode = gin.TestMode
	} else {
		Config.Mode = gin.DebugMode
	}

	if err = os.MkdirAll(path.Join(Config.Upload.Path, Config.Upload.Image.Path), PermFile); err != nil {
		return &Config, err
	}
	if err = os.MkdirAll(path.Join(Config.Upload.Path, Config.Upload.File.Path), PermFile); err != nil {
		return &Config, err
	}
	if err = os.MkdirAll(path.Join(Config.Upload.Path, Config.Upload.Image.Path, Config.Upload.Image.Thumbnail.Path), PermFile); err != nil {
		return &Config, err
	}

	InitPaths()
	InitHttp()
	InitUpload()

	return &Config, err
}
