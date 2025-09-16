package main

import (
	"fmt"
	"hs_project/models"
	"hs_project/routers"
	"os"

	"github.com/gin-gonic/gin"
)

func createLogFile() string {
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/"
	}
	if err := os.MkdirAll(logFilePath, 0777); err != nil {
		fmt.Println(err.Error())
	}
	logFileName := logFilePath + "/sql" + models.GetDate() + ".log"
	return logFileName
}

func main() {
	router := gin.Default()

	// 加载模板
	router.LoadHTMLGlob("templates/**/*")

	// 静态文件
	router.Static("/static", "./static") // 参数1：引用时使用的url，参数2：静态文件目录

	routers.AdminRoutersInit(router)

	router.Run(":8080")
}
