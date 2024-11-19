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

	start := time.Now()

	excelData, err := utils.ReadExcelData("data/data.xlsx")
	if err != nil {
		log.Fatalf("Error reading excel data: %v", err)
	}

	smtpClient, err := smtp.NewMailClient(cfg)
	if err != nil {
		log.Fatalf("Error creating mail client: %v", err)
	}

	aiClient, err := ai.NewAIClient(cfg)
	if err != nil {
		log.Fatalf("Error creating AI client: %v", err)
	}

	defer smtpClient.CloseConn()

	from := cfg.Email

	for _, row := range excelData {
		mailContent, err := aiClient.GenerateMail(row)
		if err != nil {
			log.Fatalf("Error generating mail content: %v", err)
		}

		err = smtpClient.SendMail(from, row.Email, mailContent.Subject, mailContent.HTML, "data/resume.pdf")
		if err != nil {
			log.Fatalf("Error sending mail: %v", err)
		}

		fmt.Println("Mail sent to ", row.Email, " from ", from)
		time.Sleep(time.Second * time.Duration(cfg.Delay))
	}

	elapsed := time.Since(start)
	fmt.Println("Mails sent successfully")
	fmt.Printf("Time taken: %v\n", elapsed)
}
