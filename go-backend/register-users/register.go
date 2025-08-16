package registerusers

import (
	"context"
	"echoai/db"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Creds struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Data Fetch Error: ", err)
		}

		credsjson := Creds{}

		err = json.Unmarshal(body, &credsjson)
		if err != nil {
			fmt.Println("Error Unmarshaling JSON:", err)
		}
		fmt.Println(credsjson)
		hashed, err := bcrypt.GenerateFromPassword([]byte(credsjson.Password), 12)
		if err != nil {
			fmt.Println("Error Hashing pswd:", err)
		}

		conn, err := db.Connectdb()
		if err != nil {
			fmt.Println("psql connection error from register.go:", err)
		}

		_, err = conn.Exec(context.Background(), "insert into usercreds values($1,$2,$3)", credsjson.Username, hashed, credsjson.Email)
		if err != nil {
			fmt.Println("Error Registering", err)
			fmt.Fprintf(w, "Failed Registering")
		} else {
			w.Header().Set("content-type", "text/plain")
			fmt.Fprintf(w, "success")
		}

	}
}
