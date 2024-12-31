package main

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/OM-PRAKASH-2301/ecommerce_wih_UI/handlers"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database
var tmpl *template.Template

func init() {
	// Parse the HTML templates
	tmpl = template.Must(template.ParseFiles(
		"admin/header.html",       // Include the header template
		"admin/admin_create.html", // Include your admin create page
	))
}

func main() {
	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // Your MongoDB URI here
	client, err := mongo.Connect(context.Background(), clientOptions)       // Changed nil to context.Background()
	if err != nil {
		log.Fatal(err)
	}

	// Access the e-commerce database
	db = client.Database("ecommerce")

	// Initialize handlers
	handlers.InitializeHandlers(db)

	// Set up routes
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("admin/static/"))))

	// Admin routes
	r.HandleFunc("/admin/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Render the admin_create.html template
			renderTemplate(w, "admin_create.html", nil)
		} else if r.Method == http.MethodPost {
			handlers.CreateAdmin(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}).Methods("GET", "POST")

	// Admin login route
	r.HandleFunc("/admin/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			renderTemplate(w, "admin_login.html", nil)
		} else if r.Method == http.MethodPost {
			handlers.CreateAdmin(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}).Methods("GET", "POST")

	// Other routes...
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to the E-commerce Project"))
	}).Methods("GET")

	// Start server
	http.Handle("/", r)
	log.Println("Server running at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// renderTemplate renders the templates with the given name and data
func renderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	// Execute the template (header + page)
	err := tmpl.ExecuteTemplate(w, tmplName, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
