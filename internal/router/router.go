package router

import (
	"context"
	"slices"

	"github.com/cyb3rkh4l1d/subsapi/internal/config"
	"github.com/cyb3rkh4l1d/subsapi/internal/handlers"
	"github.com/cyb3rkh4l1d/subsapi/internal/validations"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// RouteRegistrationFunc defines the function signature used to register routes.
// Функция RouteRegistrationFunc определяет сигнатуру функции, используемой для регистрации маршрутов.
type RouteRegistrationFunc func(a *Router)

// Router represents the main application container.
// It holds shared context, config, logger, HTTP engine, and handlers.
// Маршрутизатор представляет собой основной контейнер приложения.
// Он содержит общий контекст, конфигурацию, логгер, HTTP-движок и обработчики.
type Router struct {
	ctx       context.Context
	GinEngine *gin.Engine
	Logger    *logrus.Entry
	config    *config.Config
	Handler   *handlers.SubscriptionHandler
}

// NewApiRouter creates and configures the router instance.
// NewApiRouter создает и настраивает экземпляр маршрутизатора.
func NewApiRouter(ctx context.Context, config *config.Config, logger *logrus.Entry, handler *handlers.SubscriptionHandler) *Router {

	// Validate against allowed Gin modes
	// Проверка на соответствие разрешенным режимам Gin
	ginMode := ""
	ginModes := []string{gin.ReleaseMode, gin.DebugMode, gin.TestMode}

	if exists := slices.Contains(ginModes, config.GinMode); !exists {
		logger.Warnf("%+v: %+v, %+v", validations.ErrInvalidGinMode, config.GinMode, "Falling back to 'debug'.")
		ginMode = gin.DebugMode

	} else {
		ginMode = config.GinMode
	}
	gin.SetMode(ginMode)
	logger.Infof("GinMode set to : %+v", ginMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	return &Router{
		GinEngine: router,
		config:    config,
		Handler:   handler,
		Logger:    logger,
		ctx:       ctx,
	}
}

/*.....................................................................

					Functions/Methods Definations

........................................................................*/

// RegisterRoutes registers all route modules into the router instance.
// Функция RegisterRoutes регистрирует все модули маршрутизации в экземпляре маршрутизатора.
func (r *Router) RegisterRoutes(registerFuncs ...RouteRegistrationFunc) {
	for _, registerFunc := range registerFuncs {
		registerFunc(r)
	}
}
