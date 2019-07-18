package function

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/openfaas-incubator/go-function-sdk"
)

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {
	log.Println("ADD INPUT:", string(req.Body))
	values := bytes.Split(req.Body, []byte(","))
	var ret float64
	for i := 0; i < len(values); i++ {
		var value float64
		fmt.Sscanf(string(values[i]), "%f", &value)
		ret += value
	}
	var err error
	message := fmt.Sprintf("%v", ret)

	return handler.Response{
		Body:       []byte(message),
		StatusCode: http.StatusOK,
	}, err
}
