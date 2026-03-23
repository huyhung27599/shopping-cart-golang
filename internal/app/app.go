package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-management-api/internal/config"
	"user-management-api/internal/db"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/routes"
	"user-management-api/internal/validation"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

func NewApplication(cfg *config.Config) *Application {
	if err := validation.InitValidator(); err != nil {
		log.Fatalf("Validator init failed %v", err)
	}

	loadEnv()
	
	r := gin.Default()

	if err := db.InitDB(); err != nil {
		log.Fatalf("DB init failed %v", err)
	}

	redisClient := config.NewRedisClient()

	ctx := &ModuleContext{
		DB: db.DB,
		Redis: redisClient,
	}

	modules := []Module{
		NewUserModule(ctx),
	}

	routes.RegisterRoutes(r, getModulRoutes(modules)...)

	return &Application{
		config: cfg,
		router: r,
		modules: modules,
	}
}

func (a *Application) Run() error {
	srv := &http.Server{
		Addr:    a.config.ServerAddress,
		Handler: a.router,
		
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM,syscall.SIGHUP)

	go func() {
		log.Println("Starting server on", a.config.ServerAddress)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	
	}()

	<- quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}

	log.Println("Server stopped")

	return nil
}

func getModulRoutes(modules []Module) []routes.Route {
	routeList := make([]routes.Route, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}

	return routeList
}

func loadEnv() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("No .env file found")
	}
}
