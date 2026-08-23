package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/go-telegram/bot/models"
	"ye-userinfo-bot/pkg/bot"
)

// loadEnv loads environment variables from a .env file if present.
func loadEnv(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			val = strings.Trim(val, `"'`)
			if _, exists := os.LookupEnv(key); !exists {
				_ = os.Setenv(key, val)
			}
		}
	}
}

func main() {
	loadEnv(".env")

	token := strings.TrimSpace(os.Getenv("BOT_TOKEN"))
	if token == "" {
		fmt.Fprintln(os.Stderr, "Error: BOT_TOKEN environment variable is not set.")
		os.Exit(1)
	}

	b, err := bot.New(token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to initialize bot: %v\n", err)
		os.Exit(1)
	}

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "3000"
	}

	webhookPath := strings.TrimSpace(os.Getenv("WEBHOOK_PATH"))
	if webhookPath != "" && !strings.HasPrefix(webhookPath, "/") {
		webhookPath = "/" + webhookPath
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If custom secret webhook path is configured, reject any unmatched path with 404
		if webhookPath != "" && r.URL.Path != webhookPath {
			http.NotFound(w, r)
			return
		}

		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("YE UserInfo Bot is running on Vercel."))
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Verify optional secret token header from Telegram
		envSecret := strings.TrimSpace(os.Getenv("WEBHOOK_SECRET"))
		if envSecret != "" {
			incomingSecret := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
			if incomingSecret != envSecret {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var update models.Update
		if err := json.Unmarshal(body, &update); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Process update synchronously so serverless does not freeze before reply is sent
		bot.HandleUpdate(context.Background(), b, &update)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	})

	fmt.Printf("Starting HTTP webhook server on port %s (path: %s)\n", port, webhookPath)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Fprintf(os.Stderr, "HTTP server error: %v\n", err)
		os.Exit(1)
	}
}
