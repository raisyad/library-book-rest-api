package borrowing

import (
	"errors"
	"net/http"
	"strconv"

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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	filter := BorrowingFilter{
		Status: c.Query("status"),
	}

	if filter.Status != "" && filter.Status != "borrowed" && filter.Status != "returned" {
		response.Error(c, http.StatusBadRequest, "invalid status", nil)
		return
	}

	if overdue := c.Query("overdue"); overdue != "" {
		parsedOverdue, err := strconv.ParseBool(overdue)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "invalid overdue value", nil)
			return
		}
		filter.Overdue = &parsedOverdue
	}

	if memberID := c.Query("member_id"); memberID != "" {
		parsedMemberID, err := strconv.ParseInt(memberID, 10, 64)
		if err != nil || parsedMemberID <= 0 {
			response.Error(c, http.StatusBadRequest, "invalid member_id", nil)
			return
		}
		filter.MemberID = &parsedMemberID
	}

	if bookID := c.Query("book_id"); bookID != "" {
		parsedBookID, err := strconv.ParseInt(bookID, 10, 64)
		if err != nil || parsedBookID <= 0 {
			response.Error(c, http.StatusBadRequest, "invalid book_id", nil)
			return
		}
		filter.BookID = &parsedBookID
	}

	borrowings, meta, err := h.service.List(page, limit, filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to fetch borrowings", nil)
		return
	}

	response.PaginatedSuccess(c, http.StatusOK, "borrowings fetched", borrowings, meta)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := helper.ParseIDParam(c.Param("id"))
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
		response.Error(c, http.StatusBadRequest, "validation failed", validation.FormatError(err, req))
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
	id, err := helper.ParseIDParam(c.Param("id"))
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
