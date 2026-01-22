package handler

import (
	"net/http"

	"github.com/bilalabdelkadir/chis/internal/model"
	"github.com/bilalabdelkadir/chis/internal/repository"
	"github.com/bilalabdelkadir/chis/pkg/apperror"
	"github.com/bilalabdelkadir/chis/pkg/response"
	"github.com/bilalabdelkadir/chis/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo    *repository.UserRepository
	accountRepo *repository.AccountRepository
}

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName"`
}

type RegisterResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
}

func NewAuthHandler(userRepo *repository.UserRepository, accountRepo *repository.AccountRepository) *AuthHandler {
	return &AuthHandler{
		userRepo:    userRepo,
		accountRepo: accountRepo,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	var req RegisterRequest

	if err := validator.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return apperror.Internal("failed to hash password")
	}
	hashedPassword := string(hashedBytes)

	user := &model.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  &req.LastName,
	}

	if err := h.userRepo.Create(r.Context(), user); err != nil {
		return err
	}

	account := &model.Account{
		UserID:            user.ID,
		Provider:          "email",
		ProviderAccountID: user.Email,
		PasswordHash:      &hashedPassword,
	}

	if err := h.accountRepo.Create(r.Context(), account); err != nil {
		return err
	}

	res := RegisterResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
	}

	response.WriteJSON(w, http.StatusCreated, res)
	return nil
}
