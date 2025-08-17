package sessionchecker

import (
	"context"
	db "echoai/db"
	"fmt"
	"io"
	"net/http"
	"time"
)

func Check(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Err getting cookie from the frontend", err)
		}

		var sessionid string = string(body)
		fmt.Println(sessionid)
		conn, err := db.Connectdb()

		if err != nil {
			fmt.Println("Error From sessionchecker [ Error Connecting to DB ]", err)

		}
		var expiry time.Time
		conn.QueryRow(context.Background(), "select expiry_time from sessions where sessionid=$1", sessionid).Scan(expiry)
		fmt.Print(expiry)
		if expiry.After(time.Now().UTC()) {
			fmt.Fprintf(w, "Success")
		} else {
			fmt.Fprintf(w, "fail")
		}

	}

}
