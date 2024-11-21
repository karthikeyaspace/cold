package smtp

import (
	"testing"

	"github.com/karthikeyaspace/gomailer/internal/config"
)

func SmtpTest(t *testing.T) {
	cfg, err := config.LoadEnv()
	if err != nil {
		t.Fatalf("failed to load env: %v", err)
	}

	client, err := NewMailClient(cfg)
	if err != nil {
		t.Fatalf("failed to create mail client: %v", err)
	}

	err = client.SendMail(cfg.Email, cfg.Email, "Test Mail", "This is a test mail", &cfg.ResumePath)
	if err != nil {
		t.Fatalf("failed to send mail: %v", err)
	}

	t.Log("Mail sent successfully")
}
