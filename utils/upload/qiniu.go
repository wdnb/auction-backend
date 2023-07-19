package upload

import (
	"bytes"
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/spf13/viper"
)

// 上传的凭证
func getUpToken(bucket, keyToOverwrite string) string {
	accessKey := viper.GetString("qiniu.access_key")
	secretKey := viper.GetString("qiniu.secret_key")

	var scope string

	if keyToOverwrite == "" {
		scope = bucket
	} else {
		// 需要覆盖的文件名
		scope = fmt.Sprintf("%s:%s", bucket, keyToOverwrite)
	}
	//fmt.Println(scope)
	putPolicy := storage.PutPolicy{
		Scope: scope,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	return upToken
}

func getBucketManager() *storage.BucketManager {
	accessKey := viper.GetString("qiniu.access_key")
	secretKey := viper.GetString("qiniu.secret_key")
	mac := qbox.NewMac(accessKey, secretKey)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: true,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	//cfg.Region=&storage.ZoneHuabei
	bucketManager := storage.NewBucketManager(mac, &cfg)
	return bucketManager
}

// 上传到指定bucket
func UploadToQiniu(data []byte, bucket, fileName string, overwrite bool) error {
	var upToken string
	//如果keyToOverwrite为true表示覆盖上传
	if overwrite {
		upToken = getUpToken(bucket, fileName)
	} else {
		upToken = getUpToken(bucket, "")
	}

	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Region = &storage.ZoneHuadongZheJiang2
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}
	dataLen := int64(len(data))
	err := formUploader.Put(context.Background(), &ret, upToken, fileName, bytes.NewReader(data), dataLen, &putExtra)
	if err != nil {
		return err
	}
	//fmt.Println(ret.Key, ret.Hash)
	return nil
}

// 在Upload方法
//func fileMove(srcBucket, destBucket, srcKey, destKey string, overwrite bool) error {
//	err := fileMove(srcBucket, destBucket, srcKey, destKey, overwrite)
//	if err != nil {
//		return err
//	}
//	return nil
//}

// move操作移动到正式的目录里
func fileMove(srcBucket, destBucket, srcKey, destKey string, overwrite bool) error {
	//srcBucket := "if-pbl"
	//srcKey := "github.png"
	//目标空间可以和源空间相同，但是不能为跨机房的空间
	//destBucket := srcBucket
	//目标文件名可以和源文件名相同，也可以不同
	//destKey := "github-new.png"
	//如果目标文件存在，是否强制覆盖，如果不覆盖，默认返回614 file exists
	//overwrite := false
	bucketManager := getBucketManager()
	err := bucketManager.Move(srcBucket, srcKey, destBucket, destKey, overwrite)
	if err != nil {
		return err
	}
	return nil
}
