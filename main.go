package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wArrest/super_video/internal/api"
	"os"
	"path"
)

func init() {
	//命令行输出格式
	stdFormatter := &log.TextFormatter{
		DisableTimestamp: false,
		TimestampFormat:  "2006-01-02 15:04:05",
		ForceColors:      true,
		DisableColors:    false,
		FullTimestamp:    true,
	}
	//设置日志格式为格式
	log.SetFormatter(stdFormatter)
}

func main() {
	pwd, _ := os.Getwd()
	executeDir := flag.String("d", pwd, "execute directory")
	configFile := flag.String("c", "conf/super_video.yml", "config file, path in conf")
	flag.Parse()
	log.Println("Running in: ", *executeDir)
	log.Println("Using configuration: ", *configFile)
	configFullPath := path.Join(*executeDir, *configFile)
	//read config
	viper.SetConfigFile(configFullPath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	port:=viper.GetString("common.port")
	dsn := viper.GetString("dsn")
	accessPwd := viper.GetString("pwd_allocated")

	//初始化数据库
	dbPool, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Panicf("mysql connect error: ", err)
	}

	apiHandler := api.NewApiHandler(dbPool)
	apiHandler.BaseParams = api.BaseParams{
		AccessPwd: accessPwd,
	}
	r := gin.Default()
	r.POST("/transform", apiHandler.Transform)
	r.Run(port)
}
