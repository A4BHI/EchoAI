package login

import (
	"context"
	"crypto/rand"
	"echoai/db"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Checklog(username string, inputedpass string) bool {
	conn, err := db.Connectdb()
	if err != nil {
		fmt.Println("psql connection error from login.go:", err)
	}

	rows, err := conn.Query(context.Background(), "select password from usercreds where username = $1", username)
	if err != nil {
		fmt.Println("Error Executing Query:", err)
	}
	defer conn.Close(context.Background())
	defer rows.Close()
	var password string
	if rows.Next() {
		rows.Scan(&password)

		return bcrypt.CompareHashAndPassword([]byte(password), []byte(inputedpass)) == nil
	}

	return false
}
func setCookies(sessionid string, username string, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     username,
		Value:    sessionid,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Minute * 30),
	})

}
func session(username string) string {
	conn, err := db.Connectdb()
	if err != nil {
		fmt.Println("psql connection error from login.go:", err)
	}

	randomBytes := make([]byte, 32)
	rand.Read(randomBytes)
	sessionid := hex.EncodeToString(randomBytes)

	createdTime := time.Now()

	expiryTime := createdTime.Add(time.Minute * 30)
	conn.Query(context.Background(), "insert into sessions values($1,$2,$3,$4)", sessionid, username, createdTime, expiryTime)
	return sessionid
}

func Handlelogin(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		fmt.Println("Login endpoint hit")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Reading Error:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		login := Login{}
		jsonerr := json.Unmarshal(body, &login)
		if jsonerr != nil {
			fmt.Println("JSON ERROR:", jsonerr)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		fmt.Println("Received:", login.Username, login.Password)

		success := Checklog(login.Username, login.Password)

		if success {
			sessionid := session(login.Username)
			setCookies(sessionid, login.Username, w)
		} else {
			fmt.Fprint(w, "failed")
			fmt.Println("failed")
		}

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
