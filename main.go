package helloworld

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var alert map[string]interface{}

		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&alert)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println("Alert:", alert)

		// Do something with the Person struct...
		fmt.Fprintf(w, "Alert: %+v", alert)
	} else if r.Method == http.MethodGet {
		fmt.Fprint(w, "Hello World")
	}

}
