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
	start := time.Now()

	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	fmt.Println("Config loaded successfully: ", time.Since(start))

	excelData, err := utils.ReadExcelData("data/data.xlsx")
	if err != nil {
		log.Fatalf("Error reading excel data: %v", err)
	}

	fmt.Println("Excel data loaded successfully: ", time.Since(start))

	smtpClient, err := smtp.NewMailClient(cfg)
	if err != nil {
		log.Fatalf("Error creating mail client: %v", err)
	}

	fmt.Println("Mail client created successfully: ", time.Since(start))

	aiClient, err := ai.NewAIClient(cfg)
	if err != nil {
		log.Fatalf("Error creating AI client: %v", err)
	}

	fmt.Println("AI client created successfully: ", time.Since(start))

	defer smtpClient.CloseConn()

	from := cfg.Email

	for i, row := range excelData {
		fmt.Println(i, "th main started at", time.Since(start))

		mailContent, err := aiClient.GenerateMail(row)
		if err != nil {
			log.Fatalf("Error generating mail content: %v", err)
		}

		fmt.Println("mail content of ", i, "th main generated at", time.Since(start))

		err = smtpClient.SendMail(from, row.Email, mailContent.Subject, mailContent.HTML, "data/resume.pdf")
		if err != nil {
			log.Fatalf("Error sending mail: %v", err)
		}

		fmt.Println(i, "th main completed at", time.Since(start))

		fmt.Println("Waiting for ", cfg.Delay, " seconds")
		time.Sleep(time.Second * time.Duration(cfg.Delay))
	}

	elapsed := time.Since(start)
	fmt.Println("Mails sent successfully")
	fmt.Printf("Time taken: %v\n", elapsed)
}
