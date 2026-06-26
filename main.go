package main

import (
	"fish_shooting_admin_backend/configs"
	"fish_shooting_admin_backend/db/postgresql"
	"fish_shooting_admin_backend/handler"
	"fish_shooting_admin_backend/pkg/logs"
	"fish_shooting_admin_backend/pkg/redis"
	"fish_shooting_admin_backend/pkg/swagger"
	"fish_shooting_admin_backend/router"
	"fmt"
)

func main() {
	// load environment variable from .env file
	app_configs := configs.NewAppConfig()

	// log
	log_level := "info"
	logs.NewLog(log_level)

	// init postgresql database and connection pool
	pool, err := postgresql.ConnectDB()
	if err != nil {
		fmt.Println("Error connect database : ", err)
	}

	//init redis 
	_ = redis.NewRedis()

	// init go fiber framework, cors and handler configuration
	apps := router.New()

	// swagger
	swagger.NewSwagger(apps, app_configs.AppHost, app_configs.AppPort)

	// init router
	handler.NewServiceHandlers(apps, pool)

	// http server
	err = apps.Listen(fmt.Sprintf("%s:%d", app_configs.AppHost, app_configs.AppPort))
	if err != nil {
		fmt.Printf("%v", err)
	}
}
