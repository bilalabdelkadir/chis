package handler

import (
	"net/http"

	"github.com/bilalabdelkadir/chis/internal/auth"
	"github.com/bilalabdelkadir/chis/internal/model"
	"github.com/bilalabdelkadir/chis/internal/repository"
	"github.com/bilalabdelkadir/chis/pkg/apperror"
	"github.com/bilalabdelkadir/chis/pkg/response"
	"github.com/bilalabdelkadir/chis/pkg/slugify"
	"github.com/bilalabdelkadir/chis/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo         *repository.UserRepository
	accountRepo      *repository.AccountRepository
	jwtSecret        string
	organizationRepo *repository.OrganizationRepository
	membershipRepo   *repository.MembershipRepository
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

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func NewAuthHandler(
	userRepo *repository.UserRepository,
	accountRepo *repository.AccountRepository,
	organizationRepo *repository.OrganizationRepository,
	membershipRepo *repository.MembershipRepository,
	jwtSecret string,
) *AuthHandler {
	return &AuthHandler{
		userRepo:         userRepo,
		accountRepo:      accountRepo,
		organizationRepo: organizationRepo,
		membershipRepo:   membershipRepo,
		jwtSecret:        jwtSecret,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	var req LoginRequest

	if err := validator.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	userExists, err := h.userRepo.FindByEmail(r.Context(), req.Email)
	if err != nil {
		return apperror.Unauthorized("Email or password invalid.")
	}

	accountExists, err := h.accountRepo.FindByUserID(r.Context(), userExists.ID)
	if err != nil {
		return apperror.Unauthorized("Email or password invalid.")
	}

	err = bcrypt.CompareHashAndPassword([]byte(*accountExists.PasswordHash), []byte(req.Password))

	if err != nil {
		return apperror.Unauthorized("Email or password invalid.")
	}

	token, err := auth.GenerateToken(userExists.ID,
		h.jwtSecret,
	)

	if err != nil {
		return apperror.BadRequest("something error occurred please try again.")
	}

	res := LoginResponse{
		Token: token,
	}

	response.WriteJSON(w, http.StatusAccepted, res)

	return nil
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

	orgName := user.FirstName + "'s Workspace"
	slug := slugify.Slugify(orgName)

	org := &model.Organization{
		Name: orgName,
		Slug: slug,
	}

	if err := h.organizationRepo.Create(r.Context(), org); err != nil {
		return err
	}

	membership := &model.Membership{
		UserID: user.ID,
		OrgID:  org.ID,
		Role:   "admin",
	}

	if err := h.membershipRepo.Create(r.Context(), membership); err != nil {
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
