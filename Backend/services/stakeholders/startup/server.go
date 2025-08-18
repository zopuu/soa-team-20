package startup

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Mihailo84/stakeholders-service/config"
	"github.com/Mihailo84/stakeholders-service/handler"
	"github.com/Mihailo84/stakeholders-service/model"
	"github.com/Mihailo84/stakeholders-service/repository"
	"github.com/Mihailo84/stakeholders-service/service"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{config: config}
}

func (server *Server) InitializeDb() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		server.config.DBHost,
		server.config.DBUser,
		server.config.DBPass,
		server.config.DBName,
		server.config.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatal(err)
	}

	return db
}

func (server *Server) Start() {
	db := server.InitializeDb()

	userRepository := &repository.UserRepository{DatabaseConnection: db}
	userService := &service.UserService{UserRepository: userRepository}
	userHandler := handler.NewUserHandler(userService)

	router := mux.NewRouter()

	router.HandleFunc("/api/users/userById/{id}", userHandler.GetById).Methods(http.MethodGet)
	router.HandleFunc("/api/users/updateUser/{id}", userHandler.UpdateUser).Methods(http.MethodPut)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:4200",
			"http://frontend:80",    
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	handlerWithCors := c.Handler(router)

	srv := &http.Server{
		Handler:      handlerWithCors,
		Addr:         fmt.Sprintf(":%s", server.config.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Stakeholders service starting on port %s", server.config.Port)
	log.Fatal(srv.ListenAndServe())
}
