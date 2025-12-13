package routes

import (
	"github.com/cyb3rkh4l1d/subsapi/internal/app"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SwaggerRoute configures the Swagger UI endpoint (ignores the repo argument)
func SwaggerRoute(app *app.App) {

	app.GinEngine.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// SubscriptionRoutes configures the subscription-specific CRUD endpoints
func SubscriptionRoutes(app *app.App) {

	subscriptions := app.GinEngine.Group("/api/v1/subscriptions")

	subscriptions.POST("/", app.Handler.CreateSubscription)
	subscriptions.GET("/", app.Handler.ListSubscriptions)
	subscriptions.GET("/:id", app.Handler.GetSubscription)
	subscriptions.PUT("/:id", app.Handler.UpdateSubscription)
	subscriptions.DELETE("/:id", app.Handler.DeleteSubscription)

	// Stats endpoint placed under /api/v1/subscriptions/sum
	subscriptions.GET("/stats", app.Handler.SumCostHandler)
}
