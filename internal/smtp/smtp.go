package smtp

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/karthikeyaspace/gomailer/internal/config"
)

type MailClient struct {
	client *smtp.Client
	auth   smtp.Auth
}

func NewMailClient(cfg *config.Config) (*MailClient, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	auth := smtp.PlainAuth("", cfg.Email, cfg.Pass, cfg.Host)

	client, err := smtp.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to dail to %s smtp address: %v", addr, err)
	}

	tls := &tls.Config{
		ServerName: cfg.Host,
	}

	if err := client.StartTLS(tls); err != nil {
		return nil, fmt.Errorf("failed to start tls: %v", err)
	}

	if err = client.Auth(auth); err != nil {
		return nil, fmt.Errorf("authentication failed: %v", err)
	}

	return &MailClient{client: client, auth: auth}, nil
}

func (mc *MailClient) SendMail(from, to, subject, body, html string) error {

	if err := mc.client.Mail(from); err != nil {
		fmt.Println("Error while sending mail from:", err)
	}

	if err := mc.client.Rcpt(to); err != nil {
		fmt.Println("Error while sending mail to:", err)
	}

	writer, err := mc.client.Data()
	if err != nil {
		fmt.Println("Error while sending data:", err)
	}

	defer writer.Close()

	msg := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n\r\n"+
			"%s\r\n%s",
		from, to, subject, body, html,
	)
	if _, err = writer.Write([]byte(msg)); err != nil {
		return fmt.Errorf("failed to write email: %w", err)
	}

	return nil

}

func (mc *MailClient) CloseConn() error {
	return mc.client.Close()
}
