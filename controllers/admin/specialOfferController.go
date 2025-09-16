package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SpecialOfferController struct{}

func (ac SpecialOfferController) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/special_offer.html", gin.H{
		"title": "节日特价",
	})
}
