package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/bilalabdelkadir/chis/internal/email"
	"github.com/bilalabdelkadir/chis/internal/model"
	"github.com/bilalabdelkadir/chis/internal/repository"
	"github.com/bilalabdelkadir/chis/pkg/apperror"
	"github.com/bilalabdelkadir/chis/pkg/helper"
	"github.com/bilalabdelkadir/chis/pkg/response"
	"github.com/bilalabdelkadir/chis/pkg/validator"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type InvitationHandler struct {
	invitationRepo repository.InvitationRepository
	membershipRepo repository.MembershipRepository
	userRepo       repository.UserRepository
	emailService   *email.EmailService
}

func NewInvitationHandler(
	invitationRepo repository.InvitationRepository,
	membershipRepo repository.MembershipRepository,
	userRepo repository.UserRepository,
	emailService *email.EmailService,
) *InvitationHandler {
	return &InvitationHandler{
		invitationRepo: invitationRepo,
		membershipRepo: membershipRepo,
		userRepo:       userRepo,
		emailService:   emailService,
	}
}

type CreateInvitationRequest struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"omitempty,oneof=admin member"`
}

type InvitationResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	ExpiresAt string `json:"expiresAt"`
	CreatedAt string `json:"createdAt"`
}

type AcceptInvitationRequest struct {
	Token string `json:"token" validate:"required"`
}

type RespondToInvitationRequest struct {
	Action string `json:"action" validate:"required,oneof=accept decline"`
}

func (h *InvitationHandler) CreateInvitation(w http.ResponseWriter, r *http.Request) error {
	userID, err := extractUserID(r)
	if err != nil {
		return err
	}

	orgID, err := extractOrgID(r)
	if err != nil {
		return err
	}

	var req CreateInvitationRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	role := req.Role
	if role == "" {
		role = "member"
	}

	// Check if invitee is already a member
	existingUser, _ := h.userRepo.FindByEmail(r.Context(), req.Email)
	if existingUser != nil {
		_, err := h.membershipRepo.FindByUserAndOrgID(r.Context(), existingUser.ID, orgID)
		if err == nil {
			return apperror.Conflict("user is already a member of this organization")
		}
	}

	// Handle existing invitation (UNIQUE constraint on org_id, email)
	existingInv, err := h.invitationRepo.FindByOrgAndEmail(r.Context(), orgID, req.Email)
	if err == nil {
		switch existingInv.Status {
		case "pending":
			return apperror.Conflict("invitation already pending for this email")
		case "accepted":
			return apperror.Conflict("user is already a member")
		case "expired", "cancelled":
			if err := h.invitationRepo.Delete(r.Context(), existingInv.ID); err != nil {
				return apperror.Internal("failed to remove old invitation")
			}
		}
	} else if !errors.Is(err, repository.ErrNotFound) {
		return apperror.Internal("failed to check existing invitation")
	}

	token, err := helper.GenerateSecureToken()
	if err != nil {
		return apperror.Internal("failed to generate invitation token")
	}

	invitation := &model.Invitation{
		OrgID:     orgID,
		Email:     req.Email,
		Role:      role,
		Token:     token,
		Status:    "pending",
		InvitedBy: &userID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := h.invitationRepo.Create(r.Context(), invitation); err != nil {
		return apperror.Internal("failed to create invitation")
	}

	// Send invitation email
	inviterName := "A team member"
	if inviter, err := h.userRepo.FindByID(r.Context(), userID); err == nil {
		inviterName = inviter.FirstName
	}

	orgName := "your organization"
	memberships, _ := h.membershipRepo.FindAllByUserID(r.Context(), userID)
	for _, m := range memberships {
		if m.OrgID == orgID {
			orgName = m.OrgName
			break
		}
	}

	if h.emailService != nil {
		_ = h.emailService.SendInvitation(req.Email, orgName, inviterName, token)
	}

	res := InvitationResponse{
		ID:        invitation.ID.String(),
		Email:     invitation.Email,
		Role:      invitation.Role,
		Status:    invitation.Status,
		ExpiresAt: invitation.ExpiresAt.Format(time.RFC3339),
		CreatedAt: invitation.CreatedAt.Format(time.RFC3339),
	}

	response.WriteJSON(w, http.StatusCreated, res)
	return nil
}

func (h *InvitationHandler) ListInvitations(w http.ResponseWriter, r *http.Request) error {
	orgID, err := extractOrgID(r)
	if err != nil {
		return err
	}

	invitations, err := h.invitationRepo.FindByOrgID(r.Context(), orgID)
	if err != nil {
		return apperror.Internal("failed to fetch invitations")
	}

	result := make([]InvitationResponse, len(invitations))
	for i, inv := range invitations {
		result[i] = InvitationResponse{
			ID:        inv.ID.String(),
			Email:     inv.Email,
			Role:      inv.Role,
			Status:    inv.Status,
			ExpiresAt: inv.ExpiresAt.Format(time.RFC3339),
			CreatedAt: inv.CreatedAt.Format(time.RFC3339),
		}
	}

	response.WriteJSON(w, http.StatusOK, result)
	return nil
}

func (h *InvitationHandler) CancelInvitation(w http.ResponseWriter, r *http.Request) error {
	idParam := chi.URLParam(r, "id")
	invitationID, err := uuid.Parse(idParam)
	if err != nil {
		return apperror.BadRequest("invalid invitation ID")
	}

	orgID, err := extractOrgID(r)
	if err != nil {
		return err
	}

	// Verify the invitation belongs to this org
	invitations, err := h.invitationRepo.FindByOrgID(r.Context(), orgID)
	if err != nil {
		return apperror.Internal("failed to fetch invitations")
	}

	found := false
	for _, inv := range invitations {
		if inv.ID == invitationID {
			found = true
			break
		}
	}
	if !found {
		return apperror.NotFound("invitation not found")
	}

	if err := h.invitationRepo.UpdateStatus(r.Context(), invitationID, "cancelled"); err != nil {
		return apperror.Internal("failed to cancel invitation")
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"message": "invitation cancelled"})
	return nil
}

func (h *InvitationHandler) AcceptInvitation(w http.ResponseWriter, r *http.Request) error {
	userID, err := extractUserID(r)
	if err != nil {
		return err
	}

	var req AcceptInvitationRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	invitation, err := h.invitationRepo.FindByToken(r.Context(), req.Token)
	if err != nil {
		return apperror.NotFound("invitation not found")
	}

	if invitation.Status != "pending" {
		return apperror.BadRequest("invitation is no longer valid")
	}

	if time.Now().After(invitation.ExpiresAt) {
		_ = h.invitationRepo.UpdateStatus(r.Context(), invitation.ID, "expired")
		return apperror.BadRequest("invitation has expired")
	}

	// Verify the authenticated user's email matches the invitation
	user, err := h.userRepo.FindByID(r.Context(), userID)
	if err != nil {
		return apperror.Internal("failed to look up user")
	}

	if user.Email != invitation.Email {
		return apperror.Forbidden("this invitation was sent to a different email address")
	}

	// Check if already a member
	_, err = h.membershipRepo.FindByUserAndOrgID(r.Context(), userID, invitation.OrgID)
	if err == nil {
		_ = h.invitationRepo.UpdateStatus(r.Context(), invitation.ID, "accepted")
		return apperror.Conflict("you are already a member of this organization")
	}

	// Create membership
	membership := &model.Membership{
		UserID: userID,
		OrgID:  invitation.OrgID,
		Role:   invitation.Role,
	}

	if err := h.membershipRepo.Create(r.Context(), membership); err != nil {
		return apperror.Internal("failed to create membership")
	}

	if err := h.invitationRepo.AcceptInvitation(r.Context(), invitation.ID); err != nil {
		return apperror.Internal("failed to update invitation status")
	}

	// Get org details for response
	memberships, _ := h.membershipRepo.FindAllByUserID(r.Context(), userID)
	orgName := ""
	orgSlug := ""
	for _, m := range memberships {
		if m.OrgID == invitation.OrgID {
			orgName = m.OrgName
			orgSlug = m.OrgSlug
			break
		}
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "invitation accepted",
		"orgId":   invitation.OrgID.String(),
		"orgName": orgName,
		"orgSlug": orgSlug,
		"role":    invitation.Role,
	})
	return nil
}

func (h *InvitationHandler) GetPendingInvitations(w http.ResponseWriter, r *http.Request) error {
	userID, err := extractUserID(r)
	if err != nil {
		return err
	}

	user, err := h.userRepo.FindByID(r.Context(), userID)
	if err != nil {
		return apperror.Internal("failed to look up user")
	}

	pending, err := h.invitationRepo.FindPendingByEmail(r.Context(), user.Email)
	if err != nil {
		return apperror.Internal("failed to fetch pending invitations")
	}

	if pending == nil {
		pending = []repository.PendingInvitation{}
	}

	response.WriteJSON(w, http.StatusOK, pending)
	return nil
}

func (h *InvitationHandler) RespondToInvitation(w http.ResponseWriter, r *http.Request) error {
	userID, err := extractUserID(r)
	if err != nil {
		return err
	}

	idParam := chi.URLParam(r, "id")
	invitationID, err := uuid.Parse(idParam)
	if err != nil {
		return apperror.BadRequest("invalid invitation ID")
	}

	var req RespondToInvitationRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	user, err := h.userRepo.FindByID(r.Context(), userID)
	if err != nil {
		return apperror.Internal("failed to look up user")
	}

	invitation, err := h.invitationRepo.FindByID(r.Context(), invitationID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apperror.NotFound("invitation not found")
		}
		return apperror.Internal("failed to look up invitation")
	}

	if user.Email != invitation.Email {
		return apperror.Forbidden("this invitation was sent to a different email address")
	}

	if invitation.Status != "pending" {
		return apperror.BadRequest("invitation is no longer valid")
	}

	if time.Now().After(invitation.ExpiresAt) {
		_ = h.invitationRepo.UpdateStatus(r.Context(), invitation.ID, "expired")
		return apperror.BadRequest("invitation has expired")
	}

	if req.Action == "decline" {
		if err := h.invitationRepo.UpdateStatus(r.Context(), invitation.ID, "declined"); err != nil {
			return apperror.Internal("failed to decline invitation")
		}
		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "invitation declined"})
		return nil
	}

	// Accept: create membership
	_, err = h.membershipRepo.FindByUserAndOrgID(r.Context(), userID, invitation.OrgID)
	if err == nil {
		_ = h.invitationRepo.UpdateStatus(r.Context(), invitation.ID, "accepted")
		return apperror.Conflict("you are already a member of this organization")
	}

	membership := &model.Membership{
		UserID: userID,
		OrgID:  invitation.OrgID,
		Role:   invitation.Role,
	}

	if err := h.membershipRepo.Create(r.Context(), membership); err != nil {
		return apperror.Internal("failed to create membership")
	}

	if err := h.invitationRepo.AcceptInvitation(r.Context(), invitation.ID); err != nil {
		return apperror.Internal("failed to update invitation status")
	}

	memberships, _ := h.membershipRepo.FindAllByUserID(r.Context(), userID)
	orgName := ""
	orgSlug := ""
	for _, m := range memberships {
		if m.OrgID == invitation.OrgID {
			orgName = m.OrgName
			orgSlug = m.OrgSlug
			break
		}
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "invitation accepted",
		"orgId":   invitation.OrgID.String(),
		"orgName": orgName,
		"orgSlug": orgSlug,
		"role":    invitation.Role,
	})
	return nil
}
