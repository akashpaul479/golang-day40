package session

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
)

var Sessions = make(map[string]string)

// Generate session ID
func GenerateSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// Login create session
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	SessionID := GenerateSessionID()

	Sessions[SessionID] = "Akash"

	http.SetCookie(w, &http.Cookie{
		Name:  "session_id",
		Value: SessionID,
		Path:  "/",
	})

	fmt.Println(w, "Logged in. Session created!")
}

// Protected page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	cookies, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "please login first", http.StatusUnauthorized)
		return
	}
	username := Sessions[cookies.Value]
	fmt.Fprintf(w, "Welcome %s!", username)
}

// Logout destroy session
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		delete(Sessions, cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})

	fmt.Fprintln(w, "Logged out!")
}

// main func
func SessionID() {
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/home", HomeHandler)
	http.HandleFunc("/logout", LogoutHandler)

	fmt.Println("server running on port:8080")
	http.ListenAndServe(":8080", nil)
}
