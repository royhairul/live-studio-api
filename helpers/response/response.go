package response

type BaseResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewBaseResponse(message string, data any) *BaseResponse {
	if data == nil {
		return &BaseResponse{
			Message: message,
		}
	}
	return &BaseResponse{
		Message: message,
		Data:    data,
	}
}
