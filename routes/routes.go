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
	env.Router.HandleFunc("/login", controllers.Login).Methods("POST")
	env.Router.HandleFunc("/register", controllers.Register).Methods("POST")
	env.Router.HandleFunc("/verifications/{code}", controllers.Verify).Methods("GET")
	env.Router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	})
}
