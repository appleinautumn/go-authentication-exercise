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

	//  init database
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Database initialized.")
	defer db.Close()

	//  repo
	userRepository := UserRepository.NewRepository(db)
	userService := UserService.NewService(userRepository)
	userHandler := UserHandler.NewUserHandler(userService)
	authService := AuthService.NewService(userRepository)
	authHandler := AuthHandler.NewAuthHandler(authService)

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

	log.Fatal(http.ListenAndServe(fmt.Sprintf(`:%s`, os.Getenv("APP_PORT")), r))
}

func rootEndpoint(w http.ResponseWriter, r *http.Request) {
	version := os.Getenv("APP_VERSION")
	if version == "" {
		version = "0.1.0"
	}
	fmt.Fprintf(w, "go-authentication-exercise v%s", version)
}
