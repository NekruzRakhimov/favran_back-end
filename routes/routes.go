package routes

import (
	"github.com/NekruzRakhimov/favran/pkg/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

func runAllRoutes(r *gin.Engine) {
	r.GET("/ping", PingPong)

	// Endpoint for getting all main_page content {banners, sliders, categories}
	r.GET("/", controller.GetMainPageContent)

	// Endpoint for getting list of goods of defined category
	r.GET("/goods/:category_id", controller.GetCategoryGoods)

	// Endpoint for searching
	r.GET("/search", controller.Search) //search?q=some&discount=25&price_from=100&price_to=200

	// Endpoint for getting image
	r.GET("/get_image/:f1/:f2/:image_name", controller.GetImage)

	// Endpoint for adding view to the good
	r.POST("/view/:good_id", controller.AddViewToTheGood)

	r.GET("/images/:id", controller.GetGoodInfo)
}

// PingPong Проверка
func PingPong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ping": "pong"})
}

