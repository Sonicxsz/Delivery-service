package builders

import (
	"arabic/internal/handlers"
	"arabic/internal/service"
	"arabic/internal/store"
	"arabic/pkg/security/auth"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	url = "/api/v1"
)

func BuildRoutes(r *mux.Router, store *store.Store, jwtConfig *security.JWTConfig) {
	//User
	userService := service.NewUserService(store.UserRepository(), jwtConfig)
	userHandler := handlers.NewUserHandler(userService)
	r.HandleFunc(url+"/user/register", userHandler.Create()).Methods("POST")
	r.HandleFunc(url+"/user/login", userHandler.Login()).Methods("POST")

	//Tag
	tagService := service.NewTagService(store.TagRepository())
	tagHandler := handlers.NewTagHandler(tagService)
	r.HandleFunc(url+"/tag", tagHandler.Create()).Methods("POST")
	r.HandleFunc(url+"/tag/all", tagHandler.GetAll()).Methods("GET")
	r.HandleFunc(url+"/tag/{id}", tagHandler.Delete()).Methods("DELETE")

	//Category
	categoryService := service.NewCategoryService(store.CategoryRepository())
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	r.HandleFunc(url+"/category", categoryHandler.Create()).Methods("POST")
	r.HandleFunc(url+"/category/all", categoryHandler.GetAll()).Methods("GET")
	r.HandleFunc(url+"/category/{id}", categoryHandler.Delete()).Methods("DELETE")

	//Catalog
	catalogService := service.NewCatalogService(store.CatalogRepository())
	catalogHandler := handlers.NewCatalogHandler(catalogService)
	r.HandleFunc(url+"/catalog/all", catalogHandler.GetAll).Methods("GET")
	r.HandleFunc(url+"/catalog", catalogHandler.Create).Methods("POST")
	r.HandleFunc(url+"/catalog/{id}", catalogHandler.Delete).Methods("DELETE")
	r.HandleFunc(url+"/catalog", catalogHandler.Update).Methods("PATCH")
	r.HandleFunc(url+"/catalog/{id}", catalogHandler.GetById).Methods("GET")

}

func BuildProtectedRoutes(r *mux.Router, store *store.Store, jwtConfig *security.JWTConfig) {
	JWTMiddleware := security.NewJwtMiddleware(jwtConfig)

	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(JWTMiddleware.CheckJWT)

	protected.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Token is valid"))
	}).Methods("POST")
	// Будут защищенные api роуты
}
