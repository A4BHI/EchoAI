package sessionchecker

import (
	"fmt"
	"io"
	"net/http"
)

func Check(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Err getting cookie from the frontend", err)
		}

		var sessionid string = string(body)
		fmt.Println(sessionid)

	}

}
