package api_model

type Response[T any] struct {
	Data  T              `json:"data"`
	Error *ErrorResponse `json:"error"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type Operation struct {
	Success bool `json:"success"`
}
