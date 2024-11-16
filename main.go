package main

import (
	"fmt"
	"log"
	"time"

	"github.com/karthikeyaspace/gomailer/internal/config"
	"github.com/karthikeyaspace/gomailer/internal/smtp"
)

func main() {
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	client, err := smtp.NewMailClient(cfg)
	if err != nil {
		log.Fatalf("Error creating mail client: %v", err)
	}

	defer client.CloseConn()

	from := cfg.Email
	to := "karthikeyaveruturi2005@gmail.com"
	subject := "Test Mail"
	body := "Hello, This is a test mail from gomailer"
	html := "<h1>Hello, This is a test mail from gomailer</h1>"

	start := time.Now()

	for i := 1; i <= 10; i++ {
		fmt.Println("Sending mail", i)
		if err := client.SendMail(from, to, subject, body, html); err != nil {
			log.Fatalf("Error sending mail: %v", err)
		}
	}

	elapsed := time.Since(start)

	fmt.Println("Mails sent successfully")

	fmt.Printf("Time taken: %v\n", elapsed)
}
