package main

import (
	login "echoai/auth"
	register "echoai/register-users"
	session "echoai/session-cleanup"
	"log"
	"net/http"
	"time"
)

func main() {
	go func() {
		for {
			session.CleanSession()
			time.Sleep(time.Minute * 10)
		}

	}()
	http.Handle("/", http.FileServer(http.Dir("../echo")))
	http.HandleFunc("/login", login.Handlelogin)
	http.HandleFunc("/register", register.Register)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
