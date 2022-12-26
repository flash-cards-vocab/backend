package handler_interfaces

type SuccessResponse struct {
	Result interface{} `json:"result"`
}
type ErrorResponse struct {
	Message string `json:"error"`
}
