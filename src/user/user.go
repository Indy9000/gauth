package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Indy9000/gauth/src/storage"
	uuid "github.com/satori/go.uuid"
)

// Profile defines basic user profile
type Profile struct {
	UniqueUserID string
	UserName     string
	PhotoURL     string
}

// Service will handle user profiles and authentication
type Service struct {
	sessionCache     *storage.SessionCache
	sessionKeyExpiry time.Duration
	clientID         string
}

// NewService creates Service object
func NewService(sessionExpiry time.Duration, clientID string) *Service {
	return &Service{
		sessionCache:     storage.NewSessionCache(sessionExpiry),
		sessionKeyExpiry: sessionExpiry,
		clientID:         clientID,
	}
}

func (us *Service) getProfile(sessionToken string) *Profile {
	p, found := us.sessionCache.Get(sessionToken)

	if found {
		return p.(*Profile)
	}
	return nil
}

func getSessionToken(r *http.Request) (string, error) {
	c, err := r.Cookie("session_token")
	if err != nil {
		return "", err
	}
	return c.Value, nil
}

// HandleUser handles /api/user requests
func (us *Service) HandleUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet: //get user
		sessionToken, e1 := getSessionToken(r)
		if e1 != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Printf(`{"status":"failed","reason":"%s"}\n`, e1.Error())
		}
		p := us.getProfile(sessionToken)
		if p == nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println(`{"status":"unauthorized","reason":"invalid session token"}`)
			return
		}
		b, e2 := json.Marshal(*p)
		if e2 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(`{"status":"server error","reason":"profile marshalling failed"}`)
			return
		}
		fmt.Printf("returning profile:<%s>\n", string(b))
		ii, e := w.Write(b)
		if e != nil {
			fmt.Printf("Error writing. <%s>\n", e.Error())
		}
		fmt.Printf("w.Write ii:<%d> b.len:<%d>\n", ii, len(string(b)))
	case http.MethodPost: //create user

	case http.MethodPut: //update user
	case http.MethodDelete: //delete user
	}
}

// IDTokenClaims defines the reply from oauth2.googleapis.com
type IDTokenClaims struct {
	// These six fields are included in all Google ID Tokens.
	ISS string `json:"iss"` // "https://accounts.google.com",
	SUB string `json:"sub"` // "110169484474386276334",
	AZP string `json:"azp"` // "1008719970978-hb24n2dstb40o45d4feuo2ukqmcc6381.apps.googleusercontent.com",
	AUD string `json:"aud"` // "1008719970978-hb24n2dstb40o45d4feuo2ukqmcc6381.apps.googleusercontent.com",
	IAT string `json:"iat"` // "1433978353",
	EXP string `json:"exp"` // "1433981953",

	// These seven fields are only included when the user has granted the "profile" and
	// "email" OAuth scopes to the application.
	Email         string `json:"email"`          // "testuser@gmail.com",
	EmailVerified string `json:"email_verified"` //"true",
	Name          string `json:"name"`           //"Test User",
	Picture       string `json:"picture"`        //"https://lh4.googleusercontent.com/-kYgzyAWpZzJ/ABCDEFGHI/AAAJKLMNOP/tIXL9Ir44LE/s99-c/photo.jpg",
	GivenName     string `json:"given_name"`     //"Test",
	FamilyName    string `json:"family_name"`    //"User",
	Locale        string `json:"locale"`         //"en"
}

func (us *Service) validateIDToken(idtoken string) (*IDTokenClaims, error) {
	url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", idtoken)

	resp, err := http.Get(url)
	if err != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		// fmt.Println(`{"status":"unauthorised","reason":"bearer token could not be validated with https://oauth2.googleapis.com"}`)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		// w.WriteHeader(http.StatusOK)

		body, err1 := ioutil.ReadAll(resp.Body)
		if err1 != nil {
			// w.WriteHeader(http.StatusInternalServerError)
			// fmt.Println(`{"status":"unauthorised","reason":"bearer token could not be validated with https://oauth2.googleapis.com"}`)
			return nil, err1
		}
		var claims IDTokenClaims
		err2 := json.Unmarshal(body, &claims)
		if err2 != nil {
			// w.WriteHeader(http.StatusInternalServerError)
			// fmt.Println(`{"status":"unauthorised","reason":"could not unmarshal idtoken claims response"}`)
			return nil, err2
		}

		//TODO: check AUD against google clientID
		if claims.AUD != us.clientID {
			return nil, fmt.Errorf("ClientIDs don't match. Auth failed")
		}

		return &claims, nil
	}
	return nil, fmt.Errorf("Error authenticating. Google returned: %s", resp.Status)
	// w.WriteHeader(http.StatusUnauthorized)
	// fmt.Println(`{"status":"unauthorised","reason":"authentication failed with google servers."}`)

}

func createSessionKey(claims *IDTokenClaims) (string, error) {
	// Create a new random session token
	id, err := uuid.NewV4()
	if err != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		// fmt.Fprintf(w, `{"status":"failed","reason":"unable to generate session token"}`)
		return "", err
	}
	return id.String(), nil
}

// HandleUserAuth validates the id token
func (us *Service) HandleUserAuth(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost: //validate idtoken
		bearerToken := r.Header.Get("Authorization")
		if bearerToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println(`{"status":"unauthorised","reason":"no bearer token"}`)
			return
		}
		splitToken := strings.Split(bearerToken, "Bearer")
		token := strings.TrimSpace(splitToken[1])
		fmt.Printf("found token:<%s>\n", token)

		claims, err := us.validateIDToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Printf(`{"status":"unauthorised","reason":"authentication failed error:<%s>"}\n`, err.Error())
			return
		}
		fmt.Printf("SUCCESS authenticating\n")
		sessionKey, err1 := createSessionKey(claims)
		if err1 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(`{"status":"failed","reason":"unable to generate session token"}`)
		}

		p := &Profile{
			UserName:     claims.Name,
			UniqueUserID: "google-oauth2|" + claims.SUB,
			PhotoURL:     claims.Picture,
		}

		us.sessionCache.Set(sessionKey, p)
		b, e := json.Marshal(*p)
		if e != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"status":"failed","reason":"unable marshal user profile to json"}`)
		}

		// Finally, we set the client cookie for "session_token" as the session token we just generated
		// we also set an expiry time of 120 seconds, the same as the cache
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionKey,
			Expires: time.Now().UTC().Add(us.sessionKeyExpiry),
		})

		_, e1 := w.Write(b)
		if e1 != nil { //TODO: handle error properly
			fmt.Printf(`{"status":"failed","reason":"write failed. Error<%s>"}\n`, e1.Error())
		}
	}
}
