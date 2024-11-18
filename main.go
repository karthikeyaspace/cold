package main

import (
	"fmt"
	"log"
	"time"

	"github.com/karthikeyaspace/gomailer/internal/ai"
	"github.com/karthikeyaspace/gomailer/internal/config"
	"github.com/karthikeyaspace/gomailer/internal/smtp"
	"github.com/karthikeyaspace/gomailer/internal/utils"
)

func main() {
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	excelData, err := utils.ReadExcelData("data.xlsx")
	if err != nil {
		log.Fatalf("Error reading excel data: %v", err)
	}

	client, err := smtp.NewMailClient(cfg)
	if err != nil {
		log.Fatalf("Error creating mail client: %v", err)
	}

	aiClient, err := ai.NewAIClient(cfg)
	if err != nil {
		log.Fatalf("Error creating AI client: %v", err)
	}

	defer client.CloseConn()

	from := cfg.Email
	start := time.Now()

	for _, row := range excelData {
		mailContent, err := aiClient.GenerateMail(row)
		if err != nil {
			log.Fatalf("Error generating mail content: %v", err)
		}
		
	}

	elapsed := time.Since(start)

	fmt.Println("Mails sent successfully")

	fmt.Printf("Time taken: %v\n", elapsed)
}
