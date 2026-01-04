package router

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SwaggerRoute configures the Swagger UI endpoint
// SwaggerRoute настраивает конечную точку Swagger UI
func SwaggerRoute(router *Router) {

	router.GinEngine.GET("/api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Logger.Info("/api/v1/swagger: swagger api has been added")
}
