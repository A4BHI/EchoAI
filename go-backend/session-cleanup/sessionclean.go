package sessionclean

import (
	"context"
	"echoai/db"
	"fmt"
	"time"
)

func CleanSession() {
	var sessionids []string
	conn, err := db.Connectdb()
	if err != nil {
		fmt.Println("psql connection error from sessionclean.go:", err)
	}

	rows, err := conn.Query(context.Background(), "select expiry_time,sessionid from sessions")
	if err != nil {
		fmt.Println("Session Cleanup Fetch Error: ", err)
	}

	for rows.Next() {
		var expiry_time time.Time
		var sessionid string
		rows.Scan(&expiry_time, &sessionid)

		if expiry_time.Before(time.Now()) {
			sessionids = append(sessionids, sessionid)
		}
	}

	for i := range sessionids {
		conn.Exec(context.Background(), "delete from sessions where sessionid=$1", sessionids[i])
		fmt.Printf("Session %s has been deleted \n", sessionids[i])
	}

}
