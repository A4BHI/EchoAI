package sessionchecker

import (
	"context"
	db "echoai/db"
	"fmt"
	"net/http"
	"time"
)

func Check(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		cookie, err := r.Cookie("sessionid")
		if err != nil {
			fmt.Println("error from sessioncheck cant read cookie:", err)
		}

		sessionid := cookie.Value
		conn, err := db.Connectdb()
		fmt.Println(sessionid)
		if err != nil {
			fmt.Println("Error From sessionchecker [ Error Connecting to DB ]", err)

		}
		var expiry time.Time
		conn.QueryRow(context.Background(), "select expiry_time from sessions where sessionid=$1", sessionid).Scan(&expiry)

		if expiry.After(time.Now().UTC()) {
			fmt.Fprintf(w, "Success")
		} else {
			fmt.Fprint(w, "fail")
		}
	}

}
