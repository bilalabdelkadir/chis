package handler

import (
	"net/http"

	"github.com/bilalabdelkadir/chis/internal/model"
	"github.com/bilalabdelkadir/chis/internal/repository"
	"github.com/bilalabdelkadir/chis/pkg/apperror"
	"github.com/bilalabdelkadir/chis/pkg/helper"
	"github.com/bilalabdelkadir/chis/pkg/response"
	"github.com/bilalabdelkadir/chis/pkg/validator"
)

type OrganizationHandler struct {
	organizationRepo repository.OrganizationRepository
	membershipRepo   repository.MembershipRepository
}

func NewOrganizationHandler(
	organizationRepo repository.OrganizationRepository,
	membershipRepo repository.MembershipRepository,
) *OrganizationHandler {
	return &OrganizationHandler{
		organizationRepo: organizationRepo,
		membershipRepo:   membershipRepo,
	}
}

type CreateOrgRequest struct {
	Name string `json:"name" validate:"required,min=2"`
}

type OrgResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Role string `json:"role"`
}

func (h *OrganizationHandler) ListOrgs(w http.ResponseWriter, r *http.Request) error {
	userID, err := extractUserID(r)
	if err != nil {
		return err
	}

	memberships, err := h.membershipRepo.FindAllByUserID(r.Context(), userID)
	if err != nil {
		return apperror.Internal("failed to fetch organizations")
	}

	orgs := make([]OrgResponse, len(memberships))
	for i, m := range memberships {
		orgs[i] = OrgResponse{
			ID:   m.OrgID.String(),
			Name: m.OrgName,
			Slug: m.OrgSlug,
			Role: m.Role,
		}
	}

	response.WriteJSON(w, http.StatusOK, orgs)
	return nil
}

type DeleteOrgRequest struct {
	Name string `json:"name" validate:"required"`
}

func (h *OrganizationHandler) DeleteOrg(w http.ResponseWriter, r *http.Request) error {
	orgID, err := extractOrgID(r)
	if err != nil {
		return err
	}

	var req DeleteOrgRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	org, err := h.organizationRepo.FindByID(r.Context(), orgID)
	if err != nil {
		return apperror.NotFound("organization not found")
	}

	if req.Name != org.Name {
		return apperror.BadRequest("organization name does not match")
	}

	if err := h.organizationRepo.Delete(r.Context(), orgID); err != nil {
		return apperror.Internal("failed to delete organization")
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"message": "organization deleted"})
	return nil
}

func (h *OrganizationHandler) CreateOrg(w http.ResponseWriter, r *http.Request) error {
	userID, err := extractUserID(r)
	if err != nil {
		return err
	}

	var req CreateOrgRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	slug := helper.Slugify(req.Name)

	org := &model.Organization{
		Name: req.Name,
		Slug: slug,
	}

	if err := h.organizationRepo.Create(r.Context(), org); err != nil {
		return apperror.Internal("failed to create organization")
	}

	membership := &model.Membership{
		UserID: userID,
		OrgID:  org.ID,
		Role:   "admin",
	}

	if err := h.membershipRepo.Create(r.Context(), membership); err != nil {
		return apperror.Internal("failed to create membership")
	}

	res := OrgResponse{
		ID:   org.ID.String(),
		Name: org.Name,
		Slug: org.Slug,
		Role: membership.Role,
	}

	response.WriteJSON(w, http.StatusCreated, res)
	return nil
}
