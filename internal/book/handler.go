package book

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-library-rest-api/internal/response"
	"go-library-rest-api/internal/helper"
	"go-library-rest-api/internal/validation"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(c *gin.Context) {
	books, err := h.service.List()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to fetch books", nil)
		return
	}

	response.Success(c, http.StatusOK, "books fetched", books)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := helper.ParseIDParam(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid book id", nil)
		return
	}

	book, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			response.Error(c, http.StatusNotFound, "book not found", nil)
			return
		}

		response.Error(c, http.StatusInternalServerError, "failed to fetch book", nil)
		return
	}

	response.Success(c, http.StatusOK, "book fetched", book)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateBookRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "validation failed", validation.FormatError(err, req))
		return
	}

	book, err := h.service.Create(req)
	if err != nil {
		if errors.Is(err, ErrDuplicateISBN) {
			response.Error(c, http.StatusConflict, "book isbn already exists", nil)
			return
		}

		response.Error(c, http.StatusInternalServerError, "failed to create book", nil)
		return
	}

	response.Success(c, http.StatusCreated, "book created", book)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := helper.ParseIDParam(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid book id", nil)
		return
	}

	var req UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "validation failed", validation.FormatError(err, req))
		return
	}

	book, err := h.service.Update(id, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrBookNotFound):
			response.Error(c, http.StatusNotFound, "book not found", nil)
			return
		case errors.Is(err, ErrDuplicateISBN):
			response.Error(c, http.StatusConflict, "book isbn already exists", nil)
			return
		default:
			response.Error(c, http.StatusInternalServerError, "failed to update book", nil)
			return
		}
	}

	response.Success(c, http.StatusOK, "book updated", book)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := helper.ParseIDParam(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid book id", nil)
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			response.Error(c, http.StatusNotFound, "book not found", nil)
			return
		}

		response.Error(c, http.StatusInternalServerError, "failed to delete book", nil)
		return
	}

	response.Success(c, http.StatusOK, "book deleted", nil)
}