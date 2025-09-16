package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error
var Sql *gorm.DB
var DBWithLog *gorm.DB

/*
	在 Go 语言中，init() 函数会在以下时机自动执行：
	1、包初始化时执行：当程序启动时，包被导入时会自动调用 init() 函数
	2、在 main() 函数之前：init() 函数的执行顺序总是在 main() 函数之前
	3、每个包只执行一次：无论包被导入多少次，init() 函数只会执行一次

	执行顺序：
	1. 程序启动
	2. 加载并初始化所有导入的包
	3. 执行 models 包的 init() 方法（建立数据库连接）
	4. 执行 main() 函数
	5. 程序正常运行
*/
// 当引入models包时，会自动调用init()方法
func init() {
	dsn := "Enzo:Enzo_123@tcp(127.0.0.1:3306)/haircut?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info), // 开启gorm日志
	})
	if err != nil {
		panic(err)
	}
}
