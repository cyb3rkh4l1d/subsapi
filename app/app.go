package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cyb3rkh4l1d/subsapi/internal/config"
	"github.com/cyb3rkh4l1d/subsapi/internal/database"
	"github.com/cyb3rkh4l1d/subsapi/internal/handlers"
	"github.com/cyb3rkh4l1d/subsapi/internal/repository"
	"github.com/cyb3rkh4l1d/subsapi/internal/router"
	"github.com/cyb3rkh4l1d/subsapi/internal/service"
	"github.com/cyb3rkh4l1d/subsapi/internal/validations"
	"github.com/cyb3rkh4l1d/subsapi/migrations"
	"github.com/sirupsen/logrus"
)

// App encapsulates the HTTP server lifecycle management with graceful shutdown support.
// It coordinates server startup, signal handling, error propagation, and controlled termination
// Приложение инкапсулирует управление жизненным циклом HTTP-сервера с поддержкой корректного завершения работы.
// Оно координирует запуск сервера, обработку сигналов, распространение ошибок и контролируемое завершение работы.
type App struct {
	ctx             context.Context
	Server          *http.Server
	Logger          *logrus.Entry
	shutdownTimeout time.Duration
	serverErrChan   chan error
	quitChan        chan os.Signal
}

/*.....................................................................

					Functions/Methods Definations

........................................................................*/

// NewApp initializes and configures the application using the Facade design pattern.
// It orchestrates the setup of all application components in a single entry point,
// NewApp инициализирует и настраивает приложение, используя шаблон проектирования Facade.
// Он организует настройку всех компонентов приложения в единой точке входа.
func NewApp(ctx context.Context) *App {

	//initialize loggers
	//инициализация логгеров
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	configLogger := logger.WithField("component", "Config")
	dbLogger := logger.WithField("component", "Database")
	repoLogger := logger.WithField("component", "Repository")
	serviceLogger := logger.WithField("component", "Service")
	handlerLogger := logger.WithField("component", "Handler")
	routerLogger := logger.WithField("component", "Router")
	appLogger := logger.WithField("component", "App")

	// Load configuration from .env
	// Загрузка конфигурации из файла .env
	conf := config.LoadConfig(ctx, configLogger)

	//default loglevel to info
	// Уровень логирования по умолчанию: info
	logLevel, err := logrus.ParseLevel(conf.LogLevel)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logger.SetLevel(logLevel)
	appLogger.Infof("loglevel set to %+v", logLevel)

	//DATABASE: Initialize and connect to postgres database
	dbConfig := conf.DbConfig

	driver := database.NewPostgresConnection(dbConfig, dbLogger)

	//MIGRATION: Run datbase migrations
	//MIGRATION: Выполнение миграций базы данных
	migrations.PostgreSQLMigrateSubscriptions(dbLogger)

	//REPOSITORY: Initialize repository with its logger.
	//REPOSITORY: Инициализируйте репозиторий с его логгером.
	subRepo := repository.NewSubscriptionRepository(driver.Gorm_DB, repoLogger)

	//SERVICE: Initialize service with its logger.
	//SERVICE: Инициализируйте службу с её регистратором.
	subService := service.NewSubscriptionService(subRepo, serviceLogger)

	//HANDLER: Initialize handlers with its logger
	//HANDLER: Инициализируйте обработчики с помощью соответствующего логгера.
	subHandler := handlers.NewSubscriptionHandlers(ctx, handlerLogger, subService)

	//ROUTER: Initialize router with its logger
	//МАРШРУТИЗАТОР: Инициализация маршрутизатора с его логгером
	routerInstance := router.NewApiRouter(ctx, conf, routerLogger, subHandler)
	//register routes. //регистрация маршрутов
	routerInstance.RegisterRoutes(router.SubscriptionRoutes, router.SwaggerRoute)

	server := &http.Server{Addr: conf.Host, Handler: routerInstance.GinEngine}
	app := &App{
		ctx:             ctx,
		Server:          server,
		Logger:          appLogger,
		shutdownTimeout: 30 * time.Second,
		serverErrChan:   make(chan error, 1),
		quitChan:        make(chan os.Signal, 1),
	}

	return app

}

// Run starts the HTTP server and listens on the configured port.
// Команда `run` запускает HTTP-сервер и прослушивает настроенный порт.
func (a *App) Run() error {
	//defer database.ClosePgDriverConnection(a.Logger)
	// Register OS interrupt signals for graceful shutdown
	// Регистрация сигналов прерывания ОС для корректного завершения работы
	signal.Notify(a.quitChan, os.Interrupt, syscall.SIGTERM)

	// Start HTTP server in background goroutine
	// Запуск HTTP-сервера в фоновом режиме (горутина)
	go func() {
		a.Logger.Infof("starting server at :%+v", a.Server.Addr)
		if err := a.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.serverErrChan <- err
		}
		close(a.serverErrChan)
	}()

	// Wait for shutdown trigger: server error, OS signal, or context cancellation
	// Ожидание срабатывания триггера завершения работы: ошибка сервера, сигнал операционной системы или отмена контекста
	select {
	case err := <-a.serverErrChan:
		a.Logger.WithError(err).Error(validations.ErrServerStartFailed)
		return err
	case <-a.quitChan:
		a.Logger.Info("shutdown signal received.")
	case <-a.ctx.Done():
		a.Logger.Info("context cancelled signal received.")
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		a.shutdownTimeout,
	)

	defer cancel()

	// Attempt graceful shutdown within timeout
	// Попытаться корректно завершить работу программы до истечения таймаута
	if err := a.Server.Shutdown(shutdownCtx); err != nil {
		if closeErr := a.Server.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
		a.Logger.WithError(err).Info(validations.ErrShuttingServerFailed)
		return err
	}

	a.Logger.Info("server exited gracefully.")
	return nil
}
