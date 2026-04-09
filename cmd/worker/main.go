package main

import (
	"context"
	"encoding/json"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
	"user-management-api/internal/config"
	"user-management-api/internal/utils"
	"user-management-api/pkg/logger"
	"user-management-api/pkg/mail"
	"user-management-api/pkg/rabbitmq"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type Worker struct {
	rabbitMQ rabbitmq.RabbitMQService
	mailService mail.EmailProviderService
	config *config.Config
	logger *zerolog.Logger
}

func NewWorker(cfg *config.Config) *Worker {
	log := utils.NewLoggerWithPath("rabbitmq.log", "info")
	rabbitMQ, err := rabbitmq.NewRabbitMQService(utils.GetEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"), log)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create rabbitmq service")
		return nil
	}
    
	mailLogger := utils.NewLoggerWithPath("mail.log", "info")
	factory, err := mail.NewProviderFactory(mail.ProviderMailtrap)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create mail provider factory")
		return nil
	}
	mailService, err := mail.NewMailService(cfg, mailLogger, factory)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create mail service")
		return nil
	}
	return &Worker{
		rabbitMQ: rabbitMQ,
		mailService: mailService,
		config: cfg,
		logger: log,
	}
}

func (w *Worker) Start(ctx context.Context) error{
const emailQueue = "auth_mail_queue"

var email mail.Email



handler := func(body []byte) error {
	if err := json.Unmarshal(body, &email); err != nil {
		w.logger.Error().Err(err).Msg("Failed to unmarshal email")
		return err
	}
	if err := w.mailService.SendMail(ctx, &email); err != nil {
		w.logger.Error().Err(err).Msg("Failed to send email")
		return err
	}
	return nil
}

	if err := w.rabbitMQ.Consume(ctx, emailQueue, handler); err != nil {
		w.logger.Error().Err(err).Msg("Failed to consume email")
		return err
	}

	<-ctx.Done()

	return ctx.Err()
}

func (w *Worker) Stop(ctx context.Context) error{
if err := w.rabbitMQ.Close(); err != nil {
	w.logger.Error().Err(err).Msg("Failed to close rabbitmq")
	return err
}

select {
case <-ctx.Done():
	if ctx.Err() == context.DeadlineExceeded{
		w.logger.Error().Msg("Context deadline exceeded")
		return ctx.Err()
	}
default:
	return nil
}

w.logger.Info().Msg("Worker stopped")
return nil
}

func main() {
	wd := utils.MustGetWorkkingDir()

	logFile := filepath.Join(wd, "internal", "logs", "app.log")

	logger.InitLogger(logger.LoggerConfig{
		Level:      "info",
		Filename:   logFile,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
		IsDev:      utils.GetEnv("APP_ENV", "development"),
	})
	err := godotenv.Load(filepath.Join(wd, ".env"))
	if err != nil {

		logger.Log.Error().Msgf("Failed to load environment variables: %v", err)
	}

	// Initialize configuration
	cfg := config.NewConfig()

	worker := NewWorker(cfg)
	if worker == nil {
		logger.Log.Fatal().Msg("Failed to create worker")
		return
	}
	if err := worker.Start(context.Background()); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to start worker")
	
	}
  ctx, cancel := signal.NotifyContext(  context.Background(),syscall.SIGTERM, syscall.SIGHUP)
  defer cancel()

  var wg sync.WaitGroup
  wg.Add(1)
  go func() {
	defer wg.Done()
	if err := worker.Start(ctx); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to start worker")
	}
  }()

<-ctx.Done()
logger.Log.Info().Msg("Shutting down worker...")

 shutdownCtx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
 defer cancel()

 if err := worker.Stop(shutdownCtx); err != nil {
	logger.Log.Fatal().Err(err).Msg("Failed to stop worker")
 }
wg.Wait()
logger.Log.Info().Msg("Worker stopped")
return

}