package dto

type BaseResponse struct {
	Status       bool        `json:"status"`
	Code         int         `json:"code"`
	ErrorMessage []string    `json:"errorMessage"`
	Message      string      `json:"message"`
	Payload      interface{} `json:"payload"`
}

func NewBaseResponse(code int, message string, payload interface{}) *BaseResponse {
	return &BaseResponse{
		Status:  true,
		Code:    code,
		Message: message,
		Payload: payload,
	}
}

func NewErrorResponse(code int, message string, errorMessage []string) *BaseResponse {
	return &BaseResponse{
		Status:       false,
		Code:         code,
		ErrorMessage: errorMessage,
		Message:      message,
	}
}
