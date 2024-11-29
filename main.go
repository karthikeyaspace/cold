package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/karthikeyaspace/gomailer/internal/ai"
	"github.com/karthikeyaspace/gomailer/internal/config"
	"github.com/karthikeyaspace/gomailer/internal/handler"
	"github.com/karthikeyaspace/gomailer/internal/smtp"
	"github.com/karthikeyaspace/gomailer/internal/utils"
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
	router.Use(CORS)

	h := handler.NewHandler(cfg, smtpClient, aiClient)

	router.Get("/data", h.GetData)
	router.Post("/generate", h.GenerateMail)
	router.Put("/edit", h.EditMail)
	router.Post("/send", h.SendMail)

	log.Printf("Server started at %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func Loop() {
	start := time.Now()

	
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	fmt.Println("Config loaded successfully: ", time.Since(start))

	excelData, err := utils.ReadExcelData(&cfg.DataPath)
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

		err = smtpClient.SendMail(from, row.Email, mailContent.Subject, mailContent.HTML, &cfg.ResumePath)
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

func main() {
	Server()
	// Loop()
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
