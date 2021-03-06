package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Indy9000/gauth/src/storage"
	"github.com/Indy9000/gauth/src/user"
)

func main() {
	// TODO: read this from command line
	// TODO: This clientID should match with the Google Identities registered client Id. DO NOT DISCLOSE
	clientID := "16744246186-lhn5hokpb3k1i49g9pnspo4h9h3rt28f.apps.googleusercontent.com" //this client ID is a dummy
	sessionExpiry := time.Second * 60 * 15

	cache := storage.NewSessionCache(sessionExpiry)
	us := user.NewService(cache, sessionExpiry, clientID)

	http.HandleFunc("/api/v1/user", us.HandleUser)
	http.HandleFunc("/api/v1/user/auth", us.HandleUserAuth)
	http.HandleFunc("/api/v1/user/signup", us.HandleUserAuth) // TODO: this perhaps need a separate handler that creates and entry on a `user-settings` table

	http.Handle("/", http.FileServer(http.Dir("ui/")))
	fmt.Printf("Listening on 80\n")
	http.ListenAndServe(":80", nil)
}
