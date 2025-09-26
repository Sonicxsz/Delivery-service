package builders

import (
	"arabic/internal/handlers"
	"arabic/internal/service"
	"arabic/internal/store"
	"arabic/pkg/fs"
	"arabic/pkg/security/auth"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	url = "/api/v1"
)

type Builder struct {
	Router    *mux.Router
	Store     *store.Store
	JwtConfig *security.JWTConfig
	Fs        *fs.FS
}

func BuildRoutes(b *Builder) {
	//User
	userService := service.NewUserService(b.Store.UserRepository(), b.JwtConfig)
	userHandler := handlers.NewUserHandler(userService)
	b.Router.HandleFunc(url+"/user/register", userHandler.Create()).Methods("POST")
	b.Router.HandleFunc(url+"/user/login", userHandler.Login()).Methods("POST")

	//Tag
	tagService := service.NewTagService(b.Store.TagRepository())
	tagHandler := handlers.NewTagHandler(tagService)
	b.Router.HandleFunc(url+"/tag", tagHandler.Create()).Methods("POST")
	b.Router.HandleFunc(url+"/tag/all", tagHandler.GetAll()).Methods("GET")
	b.Router.HandleFunc(url+"/tag/{id}", tagHandler.Delete()).Methods("DELETE")

	//Category
	categoryService := service.NewCategoryService(b.Store.CategoryRepository())
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	b.Router.HandleFunc(url+"/category", categoryHandler.Create()).Methods("POST")
	b.Router.HandleFunc(url+"/category/all", categoryHandler.GetAll()).Methods("GET")
	b.Router.HandleFunc(url+"/category/{id}", categoryHandler.Delete()).Methods("DELETE")

	// Защищенные роуты
	JWTMiddleware := security.NewJwtMiddleware(b.JwtConfig)

	protected := b.Router.PathPrefix("/api/v1").Subrouter()
	protected.Use(JWTMiddleware.CheckJWT)

	//Catalog
	catalogService := service.NewCatalogService(b.Store.CatalogRepository())
	catalogHandler := handlers.NewCatalogHandler(catalogService)
	b.Router.HandleFunc(url+"/catalog/all", catalogHandler.GetAll(b.Fs.Image)).Methods("GET")
	b.Router.HandleFunc(url+"/catalog/{id}", catalogHandler.GetById(b.Fs.Image)).Methods("GET")

	protected.HandleFunc("/catalog", catalogHandler.Create).Methods("POST")
	protected.HandleFunc("/catalog/{id}", catalogHandler.Delete).Methods("DELETE")
	protected.HandleFunc("/catalog", catalogHandler.Update).Methods("PATCH")
	protected.HandleFunc("/catalog/add-image", catalogHandler.AddImage(b.Fs.Image)).Methods("POST")

	// User
	protected.HandleFunc("/user/profile", userHandler.Update).Methods("PATCH")
	protected.HandleFunc("/user/profile/address", userHandler.UpdateAddress).Methods("POST")
	protected.HandleFunc("/user", userHandler.Get).Methods("GET")
}

func BuildRoutesStatic(r *mux.Router, fsPath string) {
	staticPrefix := fmt.Sprintf("/%s/", fsPath)
	r.PathPrefix(staticPrefix).Handler(http.StripPrefix(staticPrefix, http.FileServer(http.Dir(fmt.Sprintf("./%s", fsPath)))))
}
