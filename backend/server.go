package main

import (
	"mini-assets/internal/handler"
	"mini-assets/internal/model"
	"mini-assets/internal/repository"
	"mini-assets/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	db := model.InitDB()
	defer db.Close()

	r := gin.Default()

	// 设置路由
	handler.SetupAssetRoutes(r, service.NewAssetService(repository.NewAssetRepository(db)))
	handler.SetupUserRoutes(r, service.NewUserService(repository.NewUserRepository(db)))

	r.Run() // 默认在8080端口启动服务
}
