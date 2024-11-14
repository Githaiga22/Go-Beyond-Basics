package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"encoding/json"
)

//user representsa simple user in our app

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// getusers returns a list of users
func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT id, name, email FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			log.Println("Error scanning user",err)
			continue
		}
		users = append(users, user)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// DB is a global database connection pool
var DB *sql.DB

func initDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	var errConn error
	DB, errConn = sql.Open("postgres", connStr)
	if errConn != nil {
		log.Fatalf("Error opening database: %v", errConn)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Cannot connect to the database: %v", err)
	}

	fmt.Println("Database connected successfully!")
}

func main() {
	initDB()
	// Set up routes and start your server here...
	http.HandleFunc("/users", getUsers)

	//start server
	fmt.Println("starting server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
