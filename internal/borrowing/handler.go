package borrowing

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-library-rest-api/internal/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(c *gin.Context) {
	borrowings, err := h.service.List()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to fetch borrowings", nil)
		return
	}

	response.Success(c, http.StatusOK, "borrowings fetched", borrowings)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := parseIDParam(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid borrowing id", nil)
		return
	}

	borrowing, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrBorrowingNotFound) {
			response.Error(c, http.StatusNotFound, "borrowing not found", nil)
			return
		}

		response.Error(c, http.StatusInternalServerError, "failed to fetch borrowing", nil)
		return
	}

	response.Success(c, http.StatusOK, "borrowing fetched", borrowing)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateBorrowingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "validation failed", err.Error())
		return
	}

	borrowing, err := h.service.Create(req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidDueDate):
			response.Error(c, http.StatusBadRequest, "invalid due date format, use YYYY-MM-DD", nil)
			return
		case errors.Is(err, ErrMemberNotFound):
			response.Error(c, http.StatusNotFound, "member not found", nil)
			return
		case errors.Is(err, ErrBookNotFound):
			response.Error(c, http.StatusNotFound, "book not found", nil)
			return
		case errors.Is(err, ErrBookOutOfStock):
			response.Error(c, http.StatusConflict, "book is out of stock", nil)
			return
		default:
			response.Error(c, http.StatusInternalServerError, "failed to create borrowing", nil)
			return
		}
	}

	response.Success(c, http.StatusCreated, "borrowing created", borrowing)
}

func (h *Handler) Return(c *gin.Context) {
	id, err := parseIDParam(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid borrowing id", nil)
		return
	}

	borrowing, err := h.service.Return(id)
	if err != nil {
		switch {
		case errors.Is(err, ErrBorrowingNotFound):
			response.Error(c, http.StatusNotFound, "borrowing not found", nil)
			return
		case errors.Is(err, ErrBorrowingAlreadyReturned):
			response.Error(c, http.StatusConflict, "borrowing already returned", nil)
			return
		default:
			response.Error(c, http.StatusInternalServerError, "failed to return book", nil)
			return
		}
	}

	response.Success(c, http.StatusOK, "book returned", borrowing)
}

func parseIDParam(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}