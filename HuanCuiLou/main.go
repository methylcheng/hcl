package main

import (
	"huancuilou/configs"
	"huancuilou/initial"
	"huancuilou/internal/article/article_controller"
	"huancuilou/internal/article/article_repository"
	"huancuilou/internal/article/article_service"
	"huancuilou/internal/homePage/activity/activity_repository"
	"huancuilou/internal/homePage/activity/activity_service"
	"huancuilou/internal/homePage/homePage_controller"
	"huancuilou/internal/user/user_controller"
	"huancuilou/internal/user/user_repository"
	"huancuilou/internal/user/user_service"
	"huancuilou/routers"
	"log"
)

func main() {
	cfg := configs.GetConfig()
	db, err := initial.Mysql(cfg.MySQL.DSN)
	if err != nil {
		log.Fatalf("初始化数据库失败：%v", err)
	}

	// 用户相关包的依赖注入
	userRepository := user_repository.NewUserRepository(db)
	userMdbRepository := user_repository.NewUserMemoryDBRepository()
	userService := user_service.NewUserService(userRepository, &cfg, userMdbRepository)
	userController := user_controller.NewUserController(userService, cfg.Jwt)

	// 文章相关包的依赖注入
	articleRepository := article_repository.NewArticleRepository(db)
	articleService := article_service.NewArticleService(articleRepository)
	articleController := article_controller.NewArticleController(articleService)

	// 活动相关包的依赖注入
	activityRepository := activity_repository.NewActivityRepository(db)
	activityService := activity_service.NewActivityService(activityRepository)
	homePageController := homePage_controller.NewHomePageController(activityService)
	sharedResourcesController := article_controller.NewSharedResourcesController(activityService)

	// 设置统一路由
	router := routers.SetupRouters(
		userController,
		articleController,
		homePageController,
		sharedResourcesController,
	)

	// 启动服务器
	if err = router.Run(":8080"); err != nil {
		log.Fatalf("启动服务器失败：%v", err)
	}
}
