package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	AuthHandler "go-authentication-exercise/auth/handler"
	AuthService "go-authentication-exercise/auth/service"
	"go-authentication-exercise/middleware"
	UserHandler "go-authentication-exercise/user/handler"
	UserRepository "go-authentication-exercise/user/repository"
	UserService "go-authentication-exercise/user/service"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	//  load env
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	// Initialize database connection
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	log.Printf("Database initialized.")

	//  repo
	userRepository := UserRepository.NewRepository(db)
	userService := UserService.NewService(userRepository)
	userHandler := UserHandler.NewUserHandler(userService)
	authService := AuthService.NewService(userRepository)
	authHandler := AuthHandler.NewAuthHandler(authService)

	// Setup router and routes
	r := setupRouter(userHandler, authHandler)

	// Start the server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}

// rootEndpoint displays the application name and version
func rootEndpoint(w http.ResponseWriter, r *http.Request) {
	version := os.Getenv("APP_VERSION")
	if version == "" {
		version = "0.1.0"
	}
	fmt.Fprintf(w, "go-authentication-exercise v%s", version)
}

// initDB initializes the database connection
func initDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// setupRouter configures all the routes for the application
func setupRouter(userHandler UserHandler.UserHandler, authHandler AuthHandler.AuthHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", rootEndpoint)

	// auth endpoints
	authRoutes := r.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/login", authHandler.Login).Methods("POST")
	authRoutes.HandleFunc("/signup", authHandler.Signup).Methods("POST")

	// user endpoints
	userRoutes := r.PathPrefix("/user").Subrouter()
	userRoutes.Use(middleware.Authenticated)
	userRoutes.HandleFunc("/list", userHandler.List).Methods("GET")

	return r
}
