package routes

import (
	"PROJECTTEST/internal/handlers"
	"PROJECTTEST/internal/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.CORSMiddleware)

	// обязательно ДО роутов:
	r.PathPrefix("/").Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
	})
	r.HandleFunc("/", handlers.HomePage).Methods("GET")
	r.HandleFunc("/api/auth/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/auth/login", handlers.LoginHandler).Methods("POST")
	r.Handle("/api/auth/me", middleware.AuthMiddleware(http.HandlerFunc(handlers.MeHandler))).Methods("GET")

	r.HandleFunc("/api/products", handlers.ListProductsHandler).Methods("GET")
	r.HandleFunc("/api/products/search", handlers.ProductsSearchHandler).Methods("GET")
	r.HandleFunc("/api/products/{id}", handlers.ProductDetailHandler).Methods("GET")
	r.HandleFunc("/api/products/category/{category}", handlers.ProductsByCategoryHandler).Methods("GET")

	r.Handle("/api/interactions/view", middleware.AuthMiddleware(http.HandlerFunc(handlers.PostInteraction))).Methods("POST")
	r.Handle("/api/interactions/like", middleware.AuthMiddleware(http.HandlerFunc(handlers.PostInteraction))).Methods("POST")
	r.Handle("/api/interactions/purchase", middleware.AuthMiddleware(http.HandlerFunc(handlers.PostInteraction))).Methods("POST")
	r.Handle("/api/user/interactions", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetUserInteractionsHandler))).Methods("GET")

	r.Handle("/api/recommendations", middleware.AuthMiddleware(http.HandlerFunc(handlers.RecommendationsHandler))).Methods("GET")

	r.Handle("/api/product/generate-100", middleware.AuthMiddleware(http.HandlerFunc(handlers.Generate100Products))).Methods("POST")
	r.Handle("/api/product/generate-1000", middleware.AuthMiddleware(http.HandlerFunc(handlers.Generate1000Products))).Methods("POST")

	return r
}
