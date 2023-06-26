package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	AuthHandler "imp/assessment/auth/handler"
	AuthService "imp/assessment/auth/service"
	"imp/assessment/middleware"
	UserHandler "imp/assessment/user/handler"
	UserRepository "imp/assessment/user/repository"
	UserService "imp/assessment/user/service"

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
	connStr := "postgres://postgres:love@localhost/imp?sslmode=disable"
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
	fmt.Fprintf(w, "imp-assessment")
}
