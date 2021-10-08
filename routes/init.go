package routes

import (
	"github.com/NekruzRakhimov/favran/pkg/controller"
	"github.com/NekruzRakhimov/favran/utils"
	"github.com/gin-gonic/gin"
	"os"
)

func InitAllRoutes() {
	r := gin.Default()

	// Исползование CORS
	r.Use(controller.CORSMiddleware())

	// Установка Logger-а
	utils.SetLogger()

	// Форматирование логов
	utils.FormatLogs(r)

	// Статус код 500, при любых panic()
	r.Use(gin.Recovery())

	// Запуск роутов
	runAllRoutes(r)

	// Запуск сервера
	_ = r.Run(":" + os.Getenv("PORT"))
}
