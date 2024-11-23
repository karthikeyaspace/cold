package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/karthikeyaspace/gomailer/internal/ai"
	"github.com/karthikeyaspace/gomailer/internal/config"
	"github.com/karthikeyaspace/gomailer/internal/smtp"
	"github.com/karthikeyaspace/gomailer/internal/utils"
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

// "GET" return the data from excel file 
func (h *Handler) GetData(w http.ResponseWriter, r *http.Request) {
	excelData, err := utils.ReadExcelData(&h.cfg.DataPath)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(excelData)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{"success": true, "data": excelData}); err != nil {
		http.Error(w, "Failed to encode data", http.StatusInternalServerError)
		return
	}
}

// "POST" take in mail in request and generate mail subject and body
func (h *Handler) GenerateMail(w http.ResponseWriter, r *http.Request) {
	var email utils.ExcelData
	if err := json.NewDecoder(r.Body).Decode(&email); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	mailContent, err := h.aiClient.GenerateMail(email)
	if err != nil {
		http.Error(w, "Failed to generate mail", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{"success": true, "data": mailContent}); err != nil {
		http.Error(w, "Failed to encode data", http.StatusInternalServerError)
		return
	}
}

// "PUT" edit the mail subject and body give id
func (h *Handler) EditMail(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Edit Mail"))
}

// "POST" send the mail given the mail subject, body and email
func (h *Handler) SendMail(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Send Mail"))
}
