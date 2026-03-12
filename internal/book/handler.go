package book

import (
	"go-library-rest-api/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(c *gin.Context) {
	response.Success(c, http.StatusOK, "books fetched successfully", []Book{})
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	response.Success(c, http.StatusOK, "book fetched successfully", gin.H{
		"id": id,
	})
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateBookRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "validation failed", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "book payload is valid", req)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "validation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "book payload is valid", gin.H{
		"id":      id,
		"payload": req,
	})
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	response.Success(c, http.StatusOK, "book route is ready", gin.H{
		"id": id,
	})
}
