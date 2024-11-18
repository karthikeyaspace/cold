package ai

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"github.com/karthikeyaspace/gomailer/internal/config"
	"github.com/karthikeyaspace/gomailer/internal/utils"
	"google.golang.org/api/option"
)

type AIClient struct {
	client *genai.Client
	model  *genai.GenerativeModel
	ctx    context.Context
}

type AIRes struct {
	Subject string
	Body    string
	HTML    string
}

func NewAIClient(cfg *config.Config) (*AIClient, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.AiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating AI client: %v", err)
	}

	model := client.GenerativeModel("gemini-1.5-pro")
	model.SetTemperature(0.3)

	return &AIClient{
		client: client,
		model:  model,
		ctx:    ctx,
	}, nil
}

// returns subject, body, html of the mail
func (a *AIClient) GenerateMail(mailData utils.ExcelData) (AIRes, error) {
	return AIRes{"hai", "hai", "hai"}, nil
}
