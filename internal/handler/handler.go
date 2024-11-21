package handler

import (
	"net/http"

	"github.com/karthikeyaspace/gomailer/internal/ai"
	"github.com/karthikeyaspace/gomailer/internal/config"
	"github.com/karthikeyaspace/gomailer/internal/smtp"
)

type Handler struct {
	cfg        *config.Config
	smtpClient *smtp.MailClient
	aiClient   *ai.AIClient
}

func NewHandler(cfg *config.Config, smtpClient *smtp.MailClient, aiClient *ai.AIClient) *Handler {
	return &Handler{
		cfg:        cfg,
		smtpClient: smtpClient,
		aiClient:   aiClient,
	}
}

func (h *Handler) CloseConn() {
	h.smtpClient.CloseConn()
}

func (h *Handler) GetData(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get Data"))
}

func (h *Handler) GenerateMail(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Generate Mail"))
}

func (h *Handler) EditMail(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Edit Mail"))
}

func (h *Handler) SendMail(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Send Mail"))
}
