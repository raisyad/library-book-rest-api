package member

type CreateMemberRequest struct {
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     *string `json:"phone"`
}

type UpdateMemberRequest struct {
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     *string `json:"phone" binding:"required"`
}