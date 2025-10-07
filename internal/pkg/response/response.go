package response

import "reflect"

type BaseResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func NewBaseResponse(message string, data any) *BaseResponse {
	if data == nil {
		return &BaseResponse{
			Message: message,
		}
	}

	if isNil(data) {
		data = makeEmptySliceOfSameType(data)
	}

	return &BaseResponse{
		Message: message,
		Data:    data,
	}
}

func isNil(data interface{}) bool {
	if data == nil {
		return false
	}
	v := reflect.ValueOf(data)
	return v.Kind() == reflect.Slice && v.IsNil()
}

func makeEmptySliceOfSameType(data interface{}) interface{} {
	t := reflect.TypeOf(data)
	return reflect.MakeSlice(t, 0, 0).Interface()
}
