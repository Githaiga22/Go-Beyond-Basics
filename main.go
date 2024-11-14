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

func loadhomepage(w http.ResponseWriter, r *http.Request) {
	// Query the database to get users
	rows, err := DB.Query("SELECT id, name, email FROM users") // retrieves all records from the users table.
	if err != nil {
		http.Error(w, "Error fetching users: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a list to hold users
	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {  //scans the values from the current row and stores them in the user struct.
			http.Error(w, "Error scanning user: "+err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user) //values are appended to the users slice once succesfully read from DB
	}

	// HTML template to render the users
	html := "<html><head><title>Users List</title></head><body>"
	html += "<h1>Users List</h1>"
	html += "<table border='1'><tr><th>ID</th><th>Name</th><th>Email</th></tr>"

	// Loop through users and display them in a table
	for _, user := range users {
		html += fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td></tr>", user.ID, user.Name, user.Email)  //function formats the user's ID, Name, and Email into an HTML table row
	}

	html += "</table></body></html>"


	w.Header().Set("Content-Type", "text/html")  //tells the browser that the response is HTML.
	w.Write([]byte(html))  //writes the HTML to the response body, which the browser renders as a webpage.
}

func main() {
	initDB()
	// Set up routes and start your server here...
	http.HandleFunc("/users", getUsers)
	http.HandleFunc("/", loadhomepage)

	//start server
	fmt.Println("starting server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
