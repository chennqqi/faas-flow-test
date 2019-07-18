package function

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		var value float64
		fmt.Sscanf(string(body), "%f", &value)
		log.Println("square:", string(body))

		w.WriteHeader(http.StatusOK)

		w.Write([]byte(fmt.Sprintf("%v", value*value)))
		return
	}

	w.WriteHeader(http.StatusNotAcceptable)
}
