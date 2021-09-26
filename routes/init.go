package routes

import (
	"favran/pkg/controller"
	"favran/utils"
	"github.com/gin-gonic/gin"
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
	_ = r.Run(utils.AppSettings.AppParams.PortRun)
}
