package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-management-api/internal/config"
	"user-management-api/internal/db"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/routes"
	"user-management-api/internal/utils"
	"user-management-api/internal/validation"
	"user-management-api/pkg/auth"
	"user-management-api/pkg/cache"
	"user-management-api/pkg/logger"
	"user-management-api/pkg/mail"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Module interface {
	Routes() routes.Route
}

type Application struct {
	config *config.Config
	router *gin.Engine
	modules []Module
}

type ModuleContext struct {
	DB sqlc.Querier
	Redis *redis.Client
}

func NewApplication(cfg *config.Config) (*Application ,error) {
	if err := validation.InitValidator(); err != nil {
		logger.Log.Error().Msgf("Validator init failed %v", err)
		return nil, err
	}

	
	
	r := gin.Default()

	if err := db.InitDB(); err != nil {
		logger.Log.Error().Msgf("DB init failed %v", err)
		return nil, err
	}

	redisClient := config.NewRedisClient()
	cacheService := cache.NewRedisCacheService(redisClient)
	tokenService := auth.NewJWTService(cacheService)
	mailLogger := utils.NewLoggerWithPath("mail.log", "info")
	factory, err := mail.NewProviderFactory(mail.ProviderMailtrap)
	if err != nil {
		mailLogger.Error().Err(err).Msg("Failed to create mail provider factory")
		return nil, err
	}
	mailService, err := mail.NewMailService(cfg, mailLogger, factory)
	if err != nil {
		mailLogger.Error().Err(err).Msg("Failed to create mail service")
		return nil, err
	}

	ctx := &ModuleContext{
		DB: db.DB,
		Redis: redisClient,
	}

	modules := []Module{
		NewUserModule(ctx),
		NewAuthModule(ctx, tokenService, cacheService, mailService),
	}

	routes.RegisterRoutes(r, tokenService, cacheService, getModulRoutes(modules)...)

	return &Application{
		config: cfg,
		router: r,
		modules: modules,
	}, nil
}

func (a *Application) Run() error {
	srv := &http.Server{
		Addr:    a.config.ServerAddress,
		Handler: a.router,
		
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM,syscall.SIGHUP)

	go func() {
		logger.Log.Info().Msgf("Starting server on %s", a.config.ServerAddress)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Log.Error().Msgf("Failed to start server: %v", err)
		}
	
	}()

	<- quit
	logger.Log.Info().Msgf("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Error().Msgf("Failed to shutdown server: %v", err)
	}

	logger.Log.Info().Msgf("Server stopped")

	return nil
}

func getModulRoutes(modules []Module) []routes.Route {
	routeList := make([]routes.Route, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}

	return routeList
}

