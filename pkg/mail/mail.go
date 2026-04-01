package mail

import (
	"context"
	"time"
	"user-management-api/internal/config"
	"user-management-api/internal/utils"
	"user-management-api/pkg/logger"

	"github.com/rs/zerolog"
)

type Email struct {
	To       []Adderss `json:"to"`
	From     Adderss   `json:"from"`
	Subject  string    `json:"subject"`
	Text     string    `json:"text"`
	Category string    `json:"category"`
}

type Adderss struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type MailConfig struct {
	ProviderConfig map[string]any 
	ProviderType   ProviderType   
	MaxRetries     int            
	Timeout        time.Duration
	Logger *zerolog.Logger
}

type MailService struct {
	config *MailConfig
	provider EmailProviderService
	logger *zerolog.Logger
}

func NewMailService(cfg *config.Config, logger *zerolog.Logger, providerFactory ProviderFactory) (EmailProviderService, error) {
	config := &MailConfig{
		ProviderConfig: cfg.MailProviderConfig,
		ProviderType: cfg.MailProviderType,
		MaxRetries: cfg.MailMaxRetries,
		Timeout: cfg.MailTimeout,
		Logger: logger,
	}
	provider, err := providerFactory.CreateProvider(config)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to create mail provider", utils.ErrCodeInternal)
	}
return &MailService{
	config: config,
	logger: logger,
	provider: provider,
}, nil
	
}

func (s *MailService) SendMail(ctx context.Context, email *Email) error {
	traceID := logger.GetTraceID(ctx)
	start := time.Now()

	

	var lastErr error
	for attemps :=1; attemps <= s.config.MaxRetries; attemps++ {
		startAttempt := time.Now()
		err := s.provider.SendMail(ctx, email)
		if err == nil {
			s.logger.Info().Str("trace_id", traceID).Dur("duration", time.Since(startAttempt)).Interface("email", email).Msg("Email sent successfully")
			return nil
		}
		lastErr = err
		 s.logger.Warn().Str("trace_id", traceID).Dur("duration", time.Since(startAttempt)).Int("attempt", attemps).Err(err).Msg("Failed to send email")

		 time.Sleep(time.Duration(attemps) * time.Second)
	}
	s.logger.Error().Str("trace_id", traceID).Dur("duration", time.Since(start)).Err(lastErr).Msg("Failed to send email after all retries")

	return utils.WrapError(lastErr, "Failed to send email after all retries", utils.ErrCodeInternal)
}