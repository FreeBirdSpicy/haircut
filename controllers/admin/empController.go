package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmpController struct{}

func (ac EmpController) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/emp.html", gin.H{
		"title": "员工列表",
	})
}
