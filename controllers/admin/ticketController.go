package admin

import (
	"fmt"
	"hs_project/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TicketController struct{}

// 列表
func (ac TicketController) Index(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	stateStr := c.DefaultQuery("state", "1")
	keyword := c.DefaultQuery("keyword", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	state, err := strconv.Atoi(stateStr)
	if err != nil || state < 1 {
		state = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	start := (page - 1) * limit

	tickets := []models.Ticket{}
	var total int64

	// 获取数据总数和分页数据
	if keyword != "" {
		models.DB.Model(&models.Ticket{}).Where("state = ? AND name LIKE ?", state, "%"+keyword+"%").Count(&total)
		models.DB.Where("state = ? AND name LIKE ?", state, "%"+keyword+"%").Limit(limit).Offset(start).Find(&tickets)
	} else {
		models.DB.Model(&models.Ticket{}).Where("state = ?", state).Count(&total)
		models.DB.Where("state = ?", state).Limit(limit).Offset(start).Find(&tickets)
	}

	c.HTML(http.StatusOK, "admin/ticket.html", gin.H{
		"title":   "项目列表",
		"data":    tickets,
		"page":    page,
		"total":   total,
		"keyword": keyword,
	})
}

func (ac TicketController) Edit(c *gin.Context) {
	var ticket models.Ticket
	c.Bind(&ticket)
	fmt.Println(ticket)
	c.Abort()

	if ticket.Id == 0 {
		// 新增
		r := models.DB.Create(&ticket)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "新增成功",
			"data": ticket,
			"r":    r.RowsAffected, // 受影响的行数
			"err":  r.Error,        // 错误信息
		})
	} else {
		// 修改
		r := models.DB.Save(&ticket)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "修改成功",
			"data": ticket,
			"r":    r.RowsAffected, // 受影响的行数
			"err":  r.Error,        // 错误信息
		})

	}
}

func (ac TicketController) Delete(c *gin.Context) {
	idStr := c.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "参数错误",
		})
		return
	} else {
		ticket := models.Ticket{Id: id}
		models.DB.First(&ticket)
		models.DB.Model(&ticket).Update("state", 0)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "删除成功",
		})
	}
}
