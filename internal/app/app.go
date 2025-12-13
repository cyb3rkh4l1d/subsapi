package app

import (
	"context"

	"github.com/cyb3rkh4l1d/subsapi/internal/database"
	"github.com/cyb3rkh4l1d/subsapi/internal/handlers"
	"github.com/cyb3rkh4l1d/subsapi/internal/models"
	"github.com/cyb3rkh4l1d/subsapi/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// RouteRegistrationFunc defines the function signature used to register routes
// into the application.
type RouteRegistrationFunc func(a *App)

// App represents the main application container.
// It holds shared context, database connection, logger, HTTP engine,
// and initialized repositories and handlers used across the service.
// it combines all the dependencies
type App struct {
	ctx       context.Context
	GinEngine *gin.Engine
	DB        *gorm.DB

	Logger *logrus.Logger
	config *Config

	Repository *repository.SubscriptionRepository
	Handler    *handlers.SubscriptionHandler
}

// NewApiApp creates and configures the main application instance.
// It initializes logging, database connection, migrations,
// repositories, handlers, and the HTTP server
func NewApiApp(ctx context.Context, config *Config) *App {

	//LOGGER: initialize logger for logging using logrus
	logger := logrus.New()
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	//DATABASE: Initialize and connect to postgres database
	dbConfig := config.dbConfig

	dbConn, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		logger.Fatalf("[-] DB error: %v", err)
	}

	//MIGRATION: Run migrations (AutoMigrate) using gorm package.
	if err := models.MigrateSubscriptions(dbConn); err != nil {
		logger.Fatalf("[-] Migration failed: %v", err)
	}

	//REPOSITORY: Initialize repository and its logger.
	repoLogger := logrus.WithField("component", "repository")
	subRepo := repository.NewSubscriptionRepository(dbConn, repoLogger)

	//HANDLER: Initialize handlers and its logger
	handlerLogger := logrus.WithField("component", "handler")
	subHandler := handlers.NewSubscriptionHandlers(ctx, handlerLogger, subRepo)

	// GIN: Initialize Gin with recovery and logger middleware (using logrus via adapter)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// APP: Finally, return the app instance.
	return &App{
		GinEngine:  router,
		config:     config,
		DB:         dbConn,
		Repository: subRepo,
		Handler:    &subHandler,
		Logger:     logger,
		ctx:        ctx,
	}
}

/*.....................................................................

					Functions/Methods Definations

........................................................................*/

// RegisterRoutes registers all route modules into the application.
// Each provided function receives the App instance and attaches routes to Gin.
func (a *App) RegisterRoutes(registerFuncs ...RouteRegistrationFunc) {
	for _, registerFunc := range registerFuncs {
		registerFunc(a)
	}
}

// Run starts the HTTP server and listens on the configured port.
func (a *App) Run() error {
	a.Logger.Infof("[+] Starting server on port %s", a.config.Port)
	srv := a.GinEngine
	return srv.Run(a.config.Port)
}
