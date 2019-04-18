package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Indy9000/gauth/src/user"
)

func main() {
	//TODO: read this from command line
	clientID := "16744246186-lhn5hokpb3k1i49g9pnspo4h9h3rt28f.apps.googleusercontent.com"

	us := user.NewService(time.Second*120, clientID)

	http.HandleFunc("/api/user", us.HandleUser)
	http.HandleFunc("/api/user/auth", us.HandleUserAuth)

	http.Handle("/", http.FileServer(http.Dir("ui/")))
	fmt.Printf("Listening on 80\n")
	http.ListenAndServe(":80", nil)
}