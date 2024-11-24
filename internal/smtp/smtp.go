package smtp

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"os"

	"github.com/karthikeyaspace/gomailer/internal/config"
)

type MailClient struct {
	client *smtp.Client
	auth   smtp.Auth
}

type MailContent struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func NewMailClient(cfg *config.Config) (*MailClient, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	auth := smtp.PlainAuth("", cfg.Email, cfg.Pass, cfg.Host)

	client, err := smtp.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to dail to %s smtp address: %v", addr, err)
	}

	// google smtp server requires tls :)
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

func (mc *MailClient) SendMail(from, to, subject, html string, resumePath *string) error {
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

	fileContent, err := os.ReadFile(*resumePath)
	if err != nil {
		return fmt.Errorf("failed to read resume file: %w", err)
	}

	encodedFile := []byte(base64.StdEncoding.EncodeToString(fileContent))

	var body bytes.Buffer
	multipartWriter := multipart.NewWriter(&body)

	headers := map[string]string{
		"From":         from,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": fmt.Sprintf("multipart/mixed; boundary=%s", multipartWriter.Boundary()),
	}

	for k, v := range headers {
		body.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	body.WriteString("\r\n")

	htmlPart, err := multipartWriter.CreatePart(map[string][]string{
		"Content-Type": {"text/html; charset=utf-8"},
	})

	if err != nil {
		return fmt.Errorf("failed to create html part: %w", err)
	}
	if _, err = htmlPart.Write([]byte(html)); err != nil {
		return fmt.Errorf("failed to write html: %w", err)
	}

	filePart, err := multipartWriter.CreatePart(map[string][]string{
		"Content-Type":              {"application/pdf"},
		"Content-Transfer-Encoding": {"base64"},
		"Content-Disposition":       {`attachment; filename="resume.pdf"`},
	})

	if err != nil {
		return fmt.Errorf("failed to create file part: %w", err)
	}

	if _, err = filePart.Write(encodedFile); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	if err = multipartWriter.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	if _, err = writer.Write([]byte(body.Bytes())); err != nil {
		return fmt.Errorf("failed to write email: %w", err)
	}

	return nil

}

func (mc *MailClient) CloseConn() error {
	return mc.client.Close()
}
