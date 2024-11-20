package ai

import (
	"context"
	"encoding/json"
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
	HTML    string
}

// creates a new AI client
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
func (ai *AIClient) GenerateMail(mailData utils.ExcelData) (AIRes, error) {
	systemPrompt := `You are a cold mail writer and generator for intern or fulltime opportunity.
		Write a professional mail with the following structure:
		1. Introduction paragraph with student background
        2. Technical skills paragraph
		3. Closing paragraph expressing interest in the opportunity

		Guidelines:
        - Keep it under 100 words
        - Be professional and courteous and cold.
        - Personalize the approach, explain how the student's background aligns with the company's work
        - Avoid overly promotional language.
        - Include relevant technical background, only include skills that are relevant to the job.
        - Make clear connection between student's interests and company/employers requirements
		- Generate the mail as if you are the student applying for the job.
		- Generate the mail as if you are directly contacting the mail recipient, dont include and placeholders or [] in the mail.

		
		Student details: 
		- Name: KARTHIKEYA VERUTURI
		- Passinate about full-stack development, Python Programming, and Machine learning.
		- Collaborative team player with a passion for innovation and growth.
		- Prioritize writing clean code.
		- Skills in Python, Typescript, React, Postgres, Fastapi, Nextjs, Prisma, AWS s3, Docker


		Mail Recipent details: 
		- Name: %s
		- Company Name: %s
		- Position: %s
		- Company Type: %s
		- Applying Position: %s
		- Additional Info: %s
		- Mutual Interests: %s
		- Reason for Contact: %s
		- Industry: %s

		Respond only with valid JSON in the following format:
		{
			"subject": "<Subject of the mail>",
			"html": "<Body of the mail in html format>",
		}
		
		Do not respond with any thing else other than the JSON format. 
		JUST JSON repsonse should start with` + "start with ```json and end with ```"

	prompt := fmt.Sprintf(systemPrompt, mailData.Name, mailData.Company, mailData.Position, mailData.CompanyType, mailData.ApplyingPosition, mailData.AdditionalInfo, mailData.MutualInterests, mailData.ReasonForContact, mailData.Industry)

	response, err := ai.model.GenerateContent(ai.ctx, genai.Text(prompt))
	if err != nil {
		return AIRes{}, fmt.Errorf("error generating mail content: %v", err)
	}

	var res AIRes
	resContent := utils.RemoveFirstAndLastLine(string(response.Candidates[0].Content.Parts[0].(genai.Text)))

	if err := json.Unmarshal([]byte(resContent), &res); err != nil {
		return AIRes{}, fmt.Errorf("error unmarshalling response: %v", err)
	}

	return res, nil
}


