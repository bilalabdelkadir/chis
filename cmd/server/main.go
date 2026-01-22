package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bilalabdelkadir/chis/internal/config"
	"github.com/bilalabdelkadir/chis/internal/database"
	"github.com/bilalabdelkadir/chis/internal/handler"
	"github.com/bilalabdelkadir/chis/internal/repository"
	"github.com/bilalabdelkadir/chis/pkg/apperror"
	"github.com/bilalabdelkadir/chis/pkg/response"
	"github.com/bilalabdelkadir/chis/pkg/validator"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := struct {
		Status string `json:"status"`
	}{
		Status: "success",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

type TestRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func TestHandler(w http.ResponseWriter, r *http.Request) error {
	var req TestRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	response.WriteJSON(w, 200, map[string]string{"status": "ok"})
	return nil
}

// Register with:

func main() {
	var err error

	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := database.Connect(cfg.DbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	fmt.Println("database connected.")

	// Create repositories
	userRepo := repository.NewUserRepository(pool)
	accountRepo := repository.NewAccountRepository(pool)

	// Create handlers
	authHandler := handler.NewAuthHandler(userRepo, accountRepo)

	// Register routes
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/test", apperror.Wrap(TestHandler))

	http.HandleFunc("/auth/register", apperror.Wrap(authHandler.Register))

	fmt.Printf("chis server starting localhost:%s\n", cfg.Port)

	err = http.ListenAndServe(":"+cfg.Port, nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
