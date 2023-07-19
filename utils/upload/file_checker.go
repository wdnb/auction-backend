package upload

import (
	"auction-website/utils"
	"errors"
	"mime/multipart"
	"net/http"
)

type CheckConfig struct {
	MaxSize   int64
	AllowExts []string
}

func CheckFile(file *multipart.FileHeader, cfg CheckConfig) (string, error) {
	// 校验文件大小
	if file.Size > cfg.MaxSize {
		return "", errors.New("文件太大")
	}

	// 打开文件检查内容
	fileHandle, err := openFile(file)
	if err != nil {
		return "", err
	}

	// 读取前几个字节检查文件头
	header := make([]byte, 261)
	_, err = fileHandle.Read(header)
	if err != nil {
		return "", err
	}
	fileType := http.DetectContentType(header)
	_, ext := utils.SplitSuffix(fileType, "/")
	isAllow := false
	allowExts := cfg.AllowExts
	for _, e := range allowExts {
		if e == ext {
			isAllow = true
			break
		}
	}
	if !isAllow {
		return "", errors.New("文件格式不支持")
	}

	return ext, nil
}
