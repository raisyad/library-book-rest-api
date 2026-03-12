package book

type CreateBookRequest struct {
	Title         string `json:"title" binding:"required"`
	Author        string `json:"author" binding:"required"`
	ISBN          string `json:"isbn" binding:"required"`
	PublishedYear *int   `json:"published_year"`
	Stock         int    `json:"stock" binding:"gte=0"`
}

type UpdateBookRequest struct {
	Title         string `json:"title" binding:"required"`
	Author        string `json:"author" binding:"required"`
	ISBN          string `json:"isbn" binding:"required"`
	PublishedYear *int   `json:"published_year"`
	Stock         int    `json:"stock" binding:"gte=0"`
}
