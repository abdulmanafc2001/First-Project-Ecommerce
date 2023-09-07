package main

import (
	"os"

	"github.com/abdulmanafc2001/First-Project-Ecommerce/database"
	_ "github.com/abdulmanafc2001/First-Project-Ecommerce/docs"
	"github.com/abdulmanafc2001/First-Project-Ecommerce/env"
	"github.com/abdulmanafc2001/First-Project-Ecommerce/routes"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	env.Loadenv()
	database.ConnectToDatabase()
	database.SyncDatabase()
}

// @title Ecommerce API
// @version 1.0
// @discription Ecommerce API in go using Gin frame work

// @host 	localhost:3000
// @BasePath /
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	router := gin.Default()

	routes.UserRoutes(router)
	routes.AdminRoutes(router)

	//add swagger in browser
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":" + port)

}
