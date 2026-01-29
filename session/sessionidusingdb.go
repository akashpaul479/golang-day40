package session

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// connect to database
func ConnectDB() {
	var err error
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/akash")
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to database...")

}

// Generate session id
func GenerateSessionId() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// Login session ID
func Loginhandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		http.Error(w, "DB not initialized", 500)
		return
	}
	SessionID := GenerateSessionId()
	expires := time.Now().Add(1 * time.Hour)

	_, err := db.Exec("INSERT INTO sessions (session_id , username , expires_at) VALUES (? , ? , ?)", SessionID, "Akash", expires)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   SessionID,
		Expires: expires,
		Path:    "/",
	})
	fmt.Fprintln(w, "Logged in!")
}

// protected page
func Homehandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "login required", http.StatusUnauthorized)
		return
	}
	var username string
	err = db.QueryRow("SELECT username FROM sessions WHERE session_id=? AND expires_at >NOW()", cookie.Value).Scan(&username)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}
	fmt.Fprintf(w, "Welcome %s!", username)
}

// Logout Session ID
func Logouthandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		db.Exec("DELETE FROM sessions WHERE session_id=?", cookie.Value)
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
func SessionIdUsingDB() {
	ConnectDB()

	http.HandleFunc("/login", Loginhandler)
	http.HandleFunc("/home", Homehandler)
	http.HandleFunc("/logout", Logouthandler)

	fmt.Println("Server running on port:8080")
	http.ListenAndServe(":8080", nil)
}
