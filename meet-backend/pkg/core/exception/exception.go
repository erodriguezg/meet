package exception

import "fmt"

type BusinessException struct {
	Code    string
	Message string
	Details map[string]string
}

// Error implements error
func (port *BusinessException) Error() string {
	return fmt.Sprintf("%s - %s - %v", port.Code, port.Message, port.Details)
}

func newBusinessException(code string, message string, details map[string]string) error {
	return &BusinessException{code, message, details}
}
