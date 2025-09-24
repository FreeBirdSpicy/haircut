package admin

import (
	"fmt"
	"hs_project/models"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type AdminController struct{}

type MonthRevenueDetil struct {
	Day   string  `json:"day"`
	Price float64 `json:"price"`
	Num   int     `json:"num"`
}

var wg sync.WaitGroup

type monthData struct {
	Month string  `json:"month"`
	Price float64 `json:"price"`
}

func (ac AdminController) Index(c *gin.Context) {
	emp_key, _ := c.Cookie("emp_key")
	emp_key_arr := strings.Split(emp_key, ",")
	username := emp_key_arr[0]

	c.HTML(http.StatusOK, "admin/index.html", gin.H{
		"title":    "首页",
		"username": username,
	})
}

// 今日营收
func (ac AdminController) TodayRevenue(c *gin.Context) {
	// 查询今日所有营收
	revenueLogs := []models.RevenueLog{}
	today_start := models.GetDate() + " 00:00:00"
	models.DB.Where("state =1 and dated >= ?", today_start).Find(&revenueLogs)

	total_num := len(revenueLogs)
	var total_price float64
	for _, v := range revenueLogs {
		price, _ := strconv.ParseFloat(v.Price, 64)
		total_price += price
	}

	// 查询所有在售票务
	tickets := []models.Ticket{}
	models.DB.Where("state = 1").Find(&tickets)

	c.HTML(http.StatusOK, "admin/today_revenue.html", gin.H{
		"title":       "今日营收",
		"total_num":   total_num,
		"total_price": total_price,
		"tickets":     tickets,
	})
}

// 添加营收记录
func (ac AdminController) RevenueLogAdd(c *gin.Context) {
	revenueLog := models.RevenueLog{}
	c.ShouldBind(&revenueLog)

	revenueLog.Plat = 1
	revenueLog.State = 1
	revenueLog.Dated = models.GetTime()

	r := models.DB.Create(&revenueLog)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "新增成功",
		"data": revenueLog,
		"r":    r.RowsAffected, // 受影响的行数
		"err":  r.Error,        // 错误信息
	})
}

// 月度营收
func (ac AdminController) MonthRevenue(c *gin.Context) {
	// 获取当前月份
	cur_month := models.CurMonth
	month := c.DefaultQuery("month", cur_month)

	//获取最近半年的月份
	monthMenu := models.GetLastHalfYear()
	for i := range monthMenu {
		monthMenu[i].On = ""
		if monthMenu[i].Value == month {
			monthMenu[i].On = "on"
		}
	}

	result := getDetail(month)

	c.HTML(http.StatusOK, "admin/month_revenue.html", gin.H{
		"title":     "月度营收",
		"monthMenu": monthMenu,
		"month":     month,
		"data":      result,
	})
}

func getDetail(month string) []MonthRevenueDetil {
	today := models.GetDate()
	start, end := models.GetMonthFirstAndLast(month)
	result := []MonthRevenueDetil{}
	sum := MonthRevenueDetil{Day: "合计", Price: 0, Num: 0}
	for i := range 31 {
		t, _ := time.Parse("2006-01-02", end)
		day := t.AddDate(0, 0, -i).Format(models.Timestamp)
		if day > today {
			continue
		}
		if day < start {
			break
		}

		revenueLogs := []models.RevenueLog{}
		models.DB.Where("state =1 and dated like ?", day+"%").Find(&revenueLogs)
		total_num := len(revenueLogs)
		sum.Num += total_num
		var total_price float64
		for _, v := range revenueLogs {
			price, _ := strconv.ParseFloat(v.Price, 64)
			total_price += price
			sum.Price += price
		}
		item := MonthRevenueDetil{Day: day, Num: total_num, Price: total_price}
		result = append(result, item)
	}
	// 将合计放在最前面
	result = append([]MonthRevenueDetil{sum}, result...)

	return result
}

// 月度营收导出
func (ac AdminController) MonthRevenueExport(c *gin.Context) {
	month := c.DefaultQuery("month", models.CurMonth)
	data := getDetail(month)

	f := excelize.NewFile()
	if err := f.Close(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  err,
		})
		return
	}

	// 创建一个工作表
	_, err := f.NewSheet("Sheet1")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  err,
		})
		return
	}

	// 写入表头
	f.SetCellValue("Sheet1", "A1", "")
	f.SetCellValue("Sheet1", "B1", "收入(元)")
	f.SetCellValue("Sheet1", "C1", "订单数量")

	for i, item := range data {
		row := i + 2 // 从第二行开始写入数据
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), item.Day)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), item.Price)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), item.Num)
	}

	// 获取桌面路径并保存文件
	userHomeDir, _ := os.UserHomeDir()
	desktopPath := filepath.Join(userHomeDir, "Desktop")
	filename := filepath.Join(desktopPath, month+"营收.xlsx")

	// 保存 Excel 文件
	if err := f.SaveAs(filename); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"msg":      "下载成功",
		"filename": filename,
	})
}

// 月度营收柱状图
func (ac AdminController) MonthRevenueChart(c *gin.Context) {
	monthMenu := models.GetLastHalfYear()

	data := make([]monthData, len(monthMenu)) // 预分配固定长度的切片

	for i, v := range monthMenu {
		fmt.Println("处理:", v.Value)
		wg.Add(1)
		go getMonthData(i, v.Value, data)
	}

	wg.Wait()

	c.JSON(http.StatusOK, gin.H{
		"code":      200,
		"msg":       "获取成功",
		"monthMenu": monthMenu,
		"data":      data,
	})
}

func getMonthData(index int, month string, data []monthData) {
	fmt.Println("开始处理:", month)
	var m = monthData{}
	models.DB.Raw("SELECT SUM(price) price FROM revenue_log WHERE state =1 and dated like ? limit 1", month+"%").Scan(&m)
	m.Month = month

	data[index] = m
	wg.Done()
}
