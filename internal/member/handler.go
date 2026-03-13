package member

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-library-rest-api/internal/helper"
	"go-library-rest-api/internal/response"
	"go-library-rest-api/internal/validation"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(c *gin.Context) {
	members, err := h.service.List()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to fetch members", nil)
		return
	}

	response.Success(c, http.StatusOK, "members fetched", members)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := helper.ParseIDParam(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid member id", nil)
		return
	}

	member, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrMemberNotFound) {
			response.Error(c, http.StatusNotFound, "member not found", nil)
			return
		}

		response.Error(c, http.StatusInternalServerError, "failed to fetch member", nil)
		return
	}

	response.Success(c, http.StatusOK, "member fetched", member)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "validation failed", validation.FormatError(err, req))
		return
	}

	member, err := h.service.Create(req)
	if err != nil {
		if errors.Is(err, ErrDuplicateEmail) {
			response.Error(c, http.StatusConflict, "member email already exists", nil)
			return
		}

		response.Error(c, http.StatusInternalServerError, "failed to create member", nil)
		return
	}

	response.Success(c, http.StatusCreated, "member created", member)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := helper.ParseIDParam(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid member id", nil)
		return
	}

	var req UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "validation failed", validation.FormatError(err, req))
		return
	}

	member, err := h.service.Update(id, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrMemberNotFound):
			response.Error(c, http.StatusNotFound, "member not found", nil)
			return
		case errors.Is(err, ErrDuplicateEmail):
			response.Error(c, http.StatusConflict, "member email already exists", nil)
			return
		default:
			response.Error(c, http.StatusInternalServerError, "failed to update member", nil)
			return
		}
	}

	response.Success(c, http.StatusOK, "member updated", member)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := helper.ParseIDParam(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid member id", nil)
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		if errors.Is(err, ErrMemberNotFound) {
			response.Error(c, http.StatusNotFound, "member not found", nil)
			return
		}

		response.Error(c, http.StatusInternalServerError, "failed to delete member", nil)
		return
	}

	response.Success(c, http.StatusOK, "member deleted", nil)
}