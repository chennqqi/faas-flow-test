package function

import (
	"fmt"
	"log"
	"net/http"

	"github.com/openfaas-incubator/go-function-sdk"
)

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {
	log.Println("ADD INPUT:", string(req.Body))
	message := fmt.Sprintf("func-ofwatchdog( %v )", string(req.Body))

	return handler.Response{
		Body:       []byte(message),
		StatusCode: http.StatusOK,
	}, nil
}
