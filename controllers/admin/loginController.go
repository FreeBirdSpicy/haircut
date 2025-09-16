package admin

import (
	"hs_project/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginController struct{}

func (lc LoginController) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/login.html", gin.H{
		"title": "登录",
	})
}

func (lc LoginController) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 检查用户名和密码是否为空
	if username == "" || password == "" {
		c.HTML(http.StatusOK, "admin/login.html", gin.H{
			"code": 200,
			"msg":  "用户名或密码不能为空",
		})
		return
	}

	// 检查用户名和密码是否正确
	emp := models.Emp{}
	models.DB.Where("username = ? and password = ?", username, password).Find(&emp)

	// 验证用户名和密码
	if emp.Id == 0 {
		// 登录失败，返回错误信息
		c.JSON(http.StatusOK, gin.H{
			"code": 401,
			"msg":  "用户名或密码错误",
		})
	} else if emp.State == 0 {
		// 登录失败，返回错误信息
		c.JSON(http.StatusOK, gin.H{
			"code": 401,
			"msg":  "该员工已被禁用",
		})
	} else {
		// 登录成功，设置会话
		emp_key := username + "," + password
		c.SetCookie("emp_key", emp_key, 86400, "/", "localhost", false, true)

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "登录成功",
		})
	}
}
