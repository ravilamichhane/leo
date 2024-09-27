package web

type SuccessMessage struct {
	Message string `json:"message"`
}

func NewSuccessMessage(message string) SuccessMessage {
	return SuccessMessage{Message: message}
}
