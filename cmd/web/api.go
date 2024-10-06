package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mznrasil/my-blogs-be/internal/middleware"
	"github.com/mznrasil/my-blogs-be/internal/services/payments"
	"github.com/mznrasil/my-blogs-be/internal/services/posts"
	"github.com/mznrasil/my-blogs-be/internal/services/sites"
	"github.com/mznrasil/my-blogs-be/internal/services/subscriptions"
	"github.com/mznrasil/my-blogs-be/internal/services/users"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()
	subRouter.Use(middleware.LoggingMiddleware)

	usersStore := users.NewStore(s.db)
	usersHandler := users.NewHandler(usersStore)
	usersHandler.RegisterRoutes(subRouter)

	sitesStore := sites.NewStore(s.db)
	sitesHandler := sites.NewHandler(sitesStore)
	sitesHandler.RegisterRoutes(subRouter)

	postsStore := posts.NewStore(s.db)
	postsHandler := posts.NewHandler(postsStore)
	postsHandler.RegisterRoutes(subRouter)

	subscriptionsStore := subscriptions.NewStore(s.db)
	subscriptionsHandler := subscriptions.NewHandler(subscriptionsStore)
	subscriptionsHandler.RegisterRoutes(subRouter)

	paymentsStore := payments.NewStore(s.db)
	paymentsHandler := payments.NewHandler(paymentsStore)
	paymentsHandler.RegisterRoutes(subRouter)

	log.Println("Server Listening on PORT", s.addr)
	http.ListenAndServe(s.addr, subRouter)
}
