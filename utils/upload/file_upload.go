package upload

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
	"mime/multipart"
)

type UploadConfig struct {
	FileName  string   //想要保存为的文件名
	MaxSize   int64    //文件最大尺寸
	OverWrite bool     //true 覆盖同名文件 false 同名文件已存在的情况下将不会覆盖
	Bucket    string   //目标bucket
	AllowExts []string //允许的扩展名
}

// 单文件上传 支持自定义文件名
func FileUpload(c *gin.Context, cfg UploadConfig) (string, error) {
	file, err := c.FormFile("file")
	if err != nil {
		return "", err
	}
	// 校验文件大小和格式
	ext, err := CheckFile(file, CheckConfig{MaxSize: cfg.MaxSize, AllowExts: cfg.AllowExts})
	if err != nil {
		return "", err
	}

	// 打开文件
	fileHandle, err := openFile(file)
	if err != nil {
		return "", err
	}

	// 读取文件内容
	var fileByte []byte
	buf := make([]byte, 1024)
	var writer bytes.Buffer
	_, err = io.CopyBuffer(&writer, fileHandle, buf)
	if err != nil {
		return "", err
	}
	fileByte = writer.Bytes()
	newFileName := fmt.Sprintf("%s.%s", cfg.FileName, ext)
	// 上传文件
	if err = UploadToQiniu(fileByte, cfg.Bucket, newFileName, cfg.OverWrite); err != nil {
		return "", err
	}
	return newFileName, nil
}

// 多文件上传 一个上传成功一个上传失败的情况不管了 定时清tmp_bucket
func MultiUpload(c *gin.Context, cfg UploadConfig) ([]string, error) {
	multipleFiles, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}
	files := multipleFiles.File["file"]
	if files == nil {
		return nil, errors.New("无法获取文件内容")
	}
	var newFileName []string
	for _, file := range files {
		// 校验文件大小和格式
		ext, err := CheckFile(file, CheckConfig{MaxSize: cfg.MaxSize, AllowExts: cfg.AllowExts})
		if err != nil {
			return nil, err
		}
		// 打开文件
		fileHandle, err := openFile(file)
		if err != nil {
			return nil, err
		}
		// 读取文件内容
		var fileByte []byte
		buf := make([]byte, 1024)
		var writer bytes.Buffer
		_, err = io.CopyBuffer(&writer, fileHandle, buf)
		if err != nil {
			return nil, err
		}
		fileByte = writer.Bytes()
		//rs := utils.RandomString(32)
		tName := fmt.Sprintf("%s.%s", cfg.FileName, ext)
		//fmt.Println(tName)
		// 上传文件
		if err = UploadToQiniu(fileByte, cfg.Bucket, tName, cfg.OverWrite); err != nil {
			return nil, err
		}
		newFileName = append(newFileName, tName)
	}
	return newFileName, nil
}

// 移动文件到正式库
func MoveFile(srcKey string, overwrite bool) error {
	srcBucket := viper.GetString("qiniu.bucket.tmp_auciton")
	destBucket := viper.GetString("qiniu.bucket.auciton")
	destKey := srcKey
	err := fileMove(srcBucket, destBucket, srcKey, destKey, overwrite)
	if err != nil {
		return err
	}
	return nil
}

func openFile(file *multipart.FileHeader) (multipart.File, error) {
	// 打开文件
	fileHandle, err := file.Open()
	if err != nil {
		return fileHandle, err
	}
	defer fileHandle.Close()
	return fileHandle, nil
}
