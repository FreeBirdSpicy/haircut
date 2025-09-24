package routers

import (
	"hs_project/controllers/admin"
	"hs_project/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 检查登录
func checkLogin(c *gin.Context) {
	// 检查登录状态
	emp_key, err := c.Cookie("emp_key")

	if err != nil {
		// 如果在 iframe 中，则用 JS 跳出 iframe
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, `
            <script>
                if (window.self !== window.top) {
                    window.top.location.href = '/login';
                } else {
                    window.location.href = '/login';
                }
            </script>
        `)
		c.Abort()
		return
	}

	emp_key_arr := strings.Split(emp_key, ",")
	username := emp_key_arr[0]
	password := emp_key_arr[1]
	emp := models.Emp{}
	models.DB.Where("username = ? and password = ? and state = 1", username, password).Find(&emp)
	if emp.Id == 0 {
		// 如果在 iframe 中，则用 JS 跳出 iframe
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, `
            <script>
                if (window.self !== window.top) {
                    window.top.location.href = '/login';
                } else {
                    window.location.href = '/login';
                }
            </script>
        `)
		c.Abort()
		return
	}
}

func AdminRoutersInit(r *gin.Engine) {
	router := r.Group("/")
	router.GET("/login", admin.LoginController{}.Index)
	router.POST("/login", admin.LoginController{}.Login)

	router.Use(checkLogin)
	{
		router.GET("/", admin.AdminController{}.Index)

		// 今日营收
		router.GET("/todayRevenue", admin.AdminController{}.TodayRevenue)
		router.POST("/revenueLogAdd", admin.AdminController{}.RevenueLogAdd)

		// 月度营收
		router.GET("/monthRevenue", admin.AdminController{}.MonthRevenue)
		router.POST("/monthRevenueExport", admin.AdminController{}.MonthRevenueExport)
		router.POST("/MonthRevenueChart", admin.AdminController{}.MonthRevenueChart)

		// 票务
		router.GET("/ticket", admin.TicketController{}.Index)
		router.POST("/ticket_edit", admin.TicketController{}.Edit)
		router.POST("/ticket_del", admin.TicketController{}.Delete)

		// 节日特价
		router.GET("/specialOffer", admin.SpecialOfferController{}.Index)

		// 员工管理
		router.GET("/emp", admin.EmpController{}.Index)
	}
}
