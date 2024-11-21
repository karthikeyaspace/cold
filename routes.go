package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/karthikeyaspace/gomailer/internal/ai"
	"github.com/karthikeyaspace/gomailer/internal/config"
	"github.com/karthikeyaspace/gomailer/internal/handler"
	"github.com/karthikeyaspace/gomailer/internal/smtp"
)

func Server() {
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	smtpClient, err := smtp.NewMailClient(cfg)
	if err != nil {
		log.Fatalf("Error creating mail client: %v", err)
	}
	defer smtpClient.CloseConn()

	aiClient, err := ai.NewAIClient(cfg)
	if err != nil {
		log.Fatalf("Error creating AI client: %v", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	h := handler.NewHandler(cfg, smtpClient, aiClient)

	router.Get("/data", h.GetData)
	router.Get("/generate", h.GenerateMail)
	router.Put("/edit", h.EditMail)
	router.Post("/send", h.SendMail)

	log.Printf("Server started at %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
