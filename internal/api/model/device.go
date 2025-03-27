package api_model

type CreateDevice struct {
	Serial string `json:"serial" binding:"required,nonempty"`
	Name   string `json:"name" binding:"required,nonempty"`
}

type UpdateDevice struct {
	Serial string `json:"serial" binding:"omitempty,nonempty"`
	Name   string `json:"name" binding:"omitempty,nonempty"`
}
