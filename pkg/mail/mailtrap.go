package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"user-management-api/internal/utils"
	"user-management-api/pkg/logger"

	"github.com/rs/zerolog"
)

type MailtrapConfig struct {
	MailSender string `json:"mail_sender"`
	NameSender string `json:"name_sender"`
	APIKey string `json:"api_key"`
	URL string `json:"url"`
}

type MailtrapProvider struct {
	client *http.Client
	config *MailtrapConfig
	logger *zerolog.Logger
}

func NewMailtrapProvider(config *MailConfig) (EmailProviderService, error) {

	mailtrapConfig ,ok:= config.ProviderConfig["mailtrap"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("mailtrap config not found")
	}

	return &MailtrapProvider{
		client: &http.Client{
			Timeout: config.Timeout,
		},
		config: &MailtrapConfig{
			MailSender: mailtrapConfig["mail_sender"].(string),
			NameSender: mailtrapConfig["name_sender"].(string),
			APIKey: mailtrapConfig["api_key"].(string),
			URL: mailtrapConfig["url"].(string),
		},
		logger: config.Logger,

	},nil
}

func (m *MailtrapProvider) SendMail(ctx context.Context, email *Email) error {
	traceID := logger.GetTraceID(ctx)
	start := time.Now()
	time.Sleep(5 * time.Second)
	email.From = Adderss{
		Email: m.config.MailSender,
		Name: m.config.NameSender,
	}
	payload,err := json.Marshal(email)
	if err != nil {
		return utils.WrapError(err, "Failed to marshal email", utils.ErrCodeInternal)
	}
 	request, err := http.NewRequestWithContext(ctx, "POST", m.config.URL, bytes.NewReader(payload))
if err != nil {
	return utils.WrapError(err, "Failed to create request", utils.ErrCodeInternal)
}
request.Header.Add("Authorization", "Bearer "+m.config.APIKey)
request.Header.Add("Content-Type", "application/json")
 response, err := m.client.Do(request)
 if err != nil {
	m.logger.Error().Str("trace_id", traceID).Dur("duration", time.Since(start)).Err(err).Msg("Failed to send request")
	return utils.WrapError(err, "Failed to send request", utils.ErrCodeInternal)
 }
 defer response.Body.Close()
 if response.StatusCode != http.StatusOK {
	body, _ := io.ReadAll(response.Body)
	m.logger.Error().Str("trace_id", traceID).Dur("duration", time.Since(start)).Int("status_code", response.StatusCode).Bytes("body", body).Msg("Failed to send request")
	return utils.WrapError(fmt.Errorf("failed to send request: %s", string(body)), "Failed to send request", utils.ErrCodeInternal)
 }
	return nil
}