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

		// 获取真实域名
		realDomain := c.Request.Host
		if forwardedHost := c.GetHeader("X-Forwarded-Host"); forwardedHost != "" {
			realDomain = forwardedHost
		}

		// HTTPS环境下，可以设置Secure为true
		isSecure := c.Request.TLS != nil

		// 清除可能存在的旧cookie
		c.SetCookie("emp_key", "", -1, "/", realDomain, false, true)
		c.SetCookie("emp_key", emp_key, 86400, "/", realDomain, isSecure, true)

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "登录成功",
		})
	}
}
