package response

type BaseResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewBaseResponse(message string, data any) *BaseResponse {
	return &BaseResponse{
		Message: message,
		Data:    data,
	}
}
