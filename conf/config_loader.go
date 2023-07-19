package conf

import (
	mdb "auction-website/database/connectors/mongodb"
	db "auction-website/database/connectors/mysql"
	"auction-website/database/connectors/redis"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

type Config struct {
	Mysql   *db.Config
	MongoDB *mdb.Config
	Redis   *rdb.Config
}

func Init() *Config {
	return &Config{
		Mysql:   Mysql(),
		MongoDB: MongoDB(),
		Redis:   Redis(),
	}
}

func Viper() (err error) {
	configPath := GetConfigPath() + string(filepath.Separator)
	// 指定配置文件路径
	viper.AddConfigPath(configPath)
	viper.SetConfigType("toml")
	viper.SetConfigName(".app")
	err = viper.ReadInConfig() // 读取配置信息

	if err != nil {
		// 读取配置信息失败
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// 监控配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config changed...")
	})
	//r := gin.Default()
	// 访问/version的返回值会随配置文件的变化而变化
	//r.GET("/version", func(c *gin.Context) {
	//	c.String(http.StatusOK, viper.GetString("version"))
	//})

	//if err := r.Run(
	//	fmt.Sprintf(":%d", viper.GetInt("port"))); err != nil {
	//	panic(err)
	//}
	return
}

func GetConfigPath() string {
	dir := getHomePath("conf")
	if strings.Contains(dir, getTmpDir()) {
		return getCurrentAbPathByCaller()
	}
	return dir
}

// 获取系统临时目录，兼容go run
func getTmpDir() string {
	dir := os.Getenv("TEMP")
	if dir == "" {
		dir = os.Getenv("TMP")
	}
	res, _ := filepath.EvalSymlinks(dir)
	return res
}

// 获取当前执行文件绝对路径
func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

// 项目目录
func getHomePath(s string) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Failed to get current file path")
	}
	// 获取项目根目录
	projectDir := filepath.Dir(filepath.Dir(filename))
	if s != "" {
		projectDir = projectDir + string(filepath.Separator) + s
	}
	return projectDir
}
