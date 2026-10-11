package httpHandler

import (
	"net/http"

	"github.com/fuadnuri/taskhub/internal/core/application/dto"
	"github.com/fuadnuri/taskhub/internal/core/application/service"
	"github.com/fuadnuri/taskhub/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrganizationHandler struct {
	service service.IOrganizationService
}

func NewOrganizationHandler(service service.IOrganizationService) *OrganizationHandler {
	return &OrganizationHandler{service: service}
}

func (h *OrganizationHandler) Create(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		RespondError(c, err)
		return
	}

	var req dto.CreateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, err := h.service.Create(c.Request.Context(), userID, req)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondCreated(c, org)
}

func (h *OrganizationHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid org id format"})
		return
	}

	org, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, org)
}

func (h *OrganizationHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	org, err := h.service.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, org)
}

func (h *OrganizationHandler) GetAll(c *gin.Context) {
	orgs, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, orgs)
}

func (h *OrganizationHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid org id format"})
		return
	}

	var req dto.UpdateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, org)
}

func (h *OrganizationHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid org id format"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, gin.H{"message": "organization deleted successfully"})
}

func (h *OrganizationHandler) InviteMember(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid org id format"})
		return
	}

	var req dto.InviteMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inv, err := h.service.InviteMember(c.Request.Context(), id, req)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondCreated(c, inv)
}

func (h *OrganizationHandler) AcceptInvitation(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		RespondError(c, err)
		return
	}

	var req dto.AcceptInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AcceptInvitation(c.Request.Context(), req.Token, userID); err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, gin.H{"message": "invitation accepted successfully"})
}

func (h *OrganizationHandler) GetMembers(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid org id format"})
		return
	}

	members, err := h.service.GetMembers(c.Request.Context(), id)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, members)
}

func (h *OrganizationHandler) RemoveMember(c *gin.Context) {
	memberID, err := uuid.Parse(c.Param("member_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid member id format"})
		return
	}

	if err := h.service.RemoveMember(c.Request.Context(), memberID); err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, gin.H{"message": "member removed successfully"})
}

func (h *OrganizationHandler) UpdateMemberRole(c *gin.Context) {
	memberID, err := uuid.Parse(c.Param("member_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid member id format"})
		return
	}

	var req dto.UpdateMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateMemberRole(c.Request.Context(), memberID, req); err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, gin.H{"message": "member role updated successfully"})
}
