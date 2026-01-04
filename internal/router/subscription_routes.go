package router

// SubscriptionRoutes configures the subscription-specific CRUD endpoints
// SubscriptionRoutes настраивает конечные точки CRUD, специфичные для каждой подписки.
func SubscriptionRoutes(router *Router) {

	subscriptions := router.GinEngine.Group("/api/v1/subscriptions")

	subscriptions.POST("/", router.Handler.CreateSubscription)
	subscriptions.GET("/", router.Handler.ListSubscriptions)
	subscriptions.GET("/:id", router.Handler.GetSubscription)
	subscriptions.PUT("/:id", router.Handler.UpdateSubscription)
	subscriptions.DELETE("/:id", router.Handler.DeleteSubscription)
	subscriptions.GET("/summary", router.Handler.GetUserSubscriptionSummary)

	router.Logger.Info("/api/vi/subscriptions: subscriptions api has been added")
}
