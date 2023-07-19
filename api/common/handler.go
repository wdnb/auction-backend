package common

import (
	"auction-website/conf"
	"auction-website/utils"
	"auction-website/utils/req"
	"auction-website/utils/resp"
	"auction-website/utils/upload"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

type Handler struct {
	Service *Service
}

func NewHandler(c *conf.Config) *Handler {
	return &Handler{
		Service: NewService(c),
	}
}

// avatarUpload godoc
// @Tags common
// @Summary 上传头像
// @Description 上传用户头像图片
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "头像图片"
// @Success 200 {object} resp.Response{data=string} "上传成功,返回上传后的文件名"
// @Failure 400 {object} resp.ErrResponse{} "请求参数错误"
// @Failure 500 {object} resp.ErrResponse{} "上传失败"
// @Router /common/file/avatar [post]
// @Security ApiKeyAuth
// @Param Access-Token header string true "JWT token" // 添加JWT头部参数
func (h *Handler) avatarUpload(c *gin.Context) {
	uid := req.GetUid(c)
	var fileName string
	//todo 检测数据库 avatar是否已设置
	//kind := "exist1"
	var overWrite bool
	//if kind == "exist" {
	//	//todo 从数据库获取头像名进行覆盖操作
	//	fileName = "exist"
	//	overWrite = true
	//} else {
	rs := utils.RandomString(32)
	sha := utils.SHA1Hash(rs)
	fileName = fmt.Sprintf("%d/avatar/%s", uid, sha)
	overWrite = true
	//}
	allowExts := []string{"jfif", "pjpeg", "jpeg", "pjp", "jpg", "png", "gif", "bmp"}
	cfg := upload.UploadConfig{
		Bucket:    viper.GetString("qiniu.bucket.tmp_auciton"),
		MaxSize:   5 * 1024 * 1024, //5m
		AllowExts: allowExts,
		FileName:  fileName,
		OverWrite: overWrite,
	}
	fName, err := upload.MultiUpload(c, cfg)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	//fmt.Println(fName)
	//upload.MoveFile(fName[0], true)
	//todo 文件名写入数据库
	resp.DataResponse(c, fName)
}

// VideoUpload godoc
// @Tags common
// @Summary 上传视频
// @Description 上传用户视频
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "视频文件"
// @Success 200 {object} resp.Response{data=string} "上传成功,返回上传后的文件名"
// @Failure 400 {object} resp.ErrResponse{} "请求参数错误"
// @Failure 500 {object} resp.ErrResponse{} "上传失败"
// @Router /common/file/video [post]
// @Security ApiKeyAuth
// @Param Access-Token header string true "JWT token" // 添加JWT头部参数
func (h *Handler) VideoUpload(c *gin.Context) {
	uid := req.GetUid(c)
	rs := utils.RandomString(32)
	sha := utils.SHA1Hash(rs)
	fileName := fmt.Sprintf("%d/video/%s", uid, sha)
	allowExts := []string{"mp4", "flv", "f4v", "webm", "m4v", "mov", "3gp", "rm", "rmvb", "ram", "wmv", "avi", "asf", "mpg", "mpeg", "dvix", "dv", "vob", "dat", "mkv", "cpk", "qt", "fli", "flc", "mod"}
	cfg := upload.UploadConfig{
		FileName:  fileName,
		Bucket:    viper.GetString("qiniu.bucket.tmp_auciton"),
		OverWrite: true,
		MaxSize:   50 * 1024 * 1024, //50m
		AllowExts: allowExts,
	}
	fName, err := upload.MultiUpload(c, cfg)
	if err != nil {
		resp.ErrorResponse(c, http.StatusBadRequest, resp.ERROR, err)
		return
	}
	resp.DataResponse(c, fName)
}
