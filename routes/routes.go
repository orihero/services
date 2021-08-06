package routes

import (
	"../controllers"
	"../env"
	"github.com/gorilla/mux"
	"net/http"
)

//----------------------ROUTES-------------------------------
//create a mux router
func CreateRouter() {
	env.Router = mux.NewRouter()
}

//initialize all routes
func InitializeRoute() {
	adminRouter := env.Router.PathPrefix("/admin/").Subrouter()
	apiRouter := env.Router.PathPrefix("/api/").Subrouter()
	//** AUTH **
	authRouter := apiRouter.PathPrefix("/auth/").Subrouter()
	authRouter.HandleFunc("/login", controllers.Login).Methods("POST")
	authRouter.HandleFunc("/register", controllers.Register).Methods("POST")
	authRouter.HandleFunc("/verifications/{code}", controllers.Verify).Methods("GET")
	//** ADMIN **
	adminRouter.HandleFunc("/user", controllers.Register).Methods("POST")
	adminRouter.HandleFunc("/user", controllers.GetUsers).Methods("GET")
	adminRouter.HandleFunc("/user", controllers.UpdateUser).Methods("PUT")
	adminRouter.HandleFunc("/user/{id}", controllers.DeleteUser).Methods("DELETE")

	env.Router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	})
}
