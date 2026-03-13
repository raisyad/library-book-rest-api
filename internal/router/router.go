package router

import (
	"go-library-rest-api/internal/book"
	"go-library-rest-api/internal/response"
	"go-library-rest-api/internal/member"
	"go-library-rest-api/internal/borrowing"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func Setup(db *sqlx.DB) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	api := r.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/health", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			response.Error(c, http.StatusServiceUnavailable, "database is not reachable", nil)
			return
		}

		response.Success(c, http.StatusOK, "server is running", gin.H{
			"database": "connected",
		})
	})

	bookRepo := book.NewRepository(db)
	bookService := book.NewService(bookRepo)
	bookHandler := book.NewHandler(bookService)

	book.RegisterRoutes(v1, bookHandler)

	memberRepo := member.NewRepository(db)
	memberService := member.NewService(memberRepo)
	memberHandler := member.NewHandler(memberService)

	member.RegisterRoutes(v1, memberHandler)

	borrowingRepo := borrowing.NewRepository(db)
	borrowingService := borrowing.NewService(borrowingRepo)
	borrowingHandler := borrowing.NewHandler(borrowingService)

	borrowing.RegisterRoutes(v1, borrowingHandler)

	return r
}
