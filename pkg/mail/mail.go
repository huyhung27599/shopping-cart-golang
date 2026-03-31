package mail

import (
	"time"
	"user-management-api/internal/config"

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

func NewMailService(cfg *config.Config, logger *zerolog.Logger, providerFactory ProviderFactory) (EmailProviderService, error) {
	
	
}