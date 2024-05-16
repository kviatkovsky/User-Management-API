package api

import (
	"main/services/users"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type APIServer struct {
	addr  string
	db    *gorm.DB
	redis *redis.Client
}

func NewAPIServer(addr string, db *gorm.DB, redis *redis.Client) *APIServer {
	return &APIServer{
		addr:  addr,
		db:    db,
		redis: redis,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1/").Subrouter()
	userStore := users.NewStore(s.db, s.redis)
	userHelper := users.NewHelper(s.db, userStore)
	userHandler := users.NewHandler(userStore, userHelper)
	userHandler.RegisterRoutes(subRouter)

	return http.ListenAndServe(s.addr, router)
}
