package response

type meta struct {
	Message string `json:"message"`
}

type successResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type errorResponse struct {
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}
