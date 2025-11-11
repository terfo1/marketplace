package routes

import (
	"PROJECTTEST/internal/handlers"
	"PROJECTTEST/internal/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func InitRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.HomePage).Methods("GET")
	r.HandleFunc("/api/auth/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/auth/login", handlers.LoginHandler).Methods("POST")

	r.HandleFunc("/api/products", handlers.ListProductsHandler).Methods("GET")
	r.HandleFunc("/api/products/{id}", handlers.ProductDetailHandler).Methods("GET")
	r.HandleFunc("/api/products/category/{category}", handlers.ProductsByCategoryHandler).Methods("GET")
	r.HandleFunc("/api/products/search", handlers.ProductsSearchHandler).Methods("GET")

	r.Handle("/api/interactions/view", middleware.AuthMiddleware(http.HandlerFunc(handlers.PostInteraction))).Methods("POST")
	r.Handle("/api/interactions/like", middleware.AuthMiddleware(http.HandlerFunc(handlers.PostInteraction))).Methods("POST")
	r.Handle("/api/interactions/purchase", middleware.AuthMiddleware(http.HandlerFunc(handlers.PostInteraction))).Methods("POST")
	r.Handle("/api/user/interactions", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetUserInteractionsHandler))).Methods("GET")

	r.Handle("/api/recommendations", middleware.AuthMiddleware(http.HandlerFunc(handlers.RecommendationsHandler))).Methods("GET")
	return r
}
