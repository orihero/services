package main

import (
	"./database"
	"./env"
	"./routes"
	"fmt"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
)

//start the server
func ServerStart() {
	fmt.Println("Server started at http://localhost:8080")
	err := http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(env.Router))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	database.InitialMigration()
	routes.CreateRouter()
	routes.InitializeRoute()
	ServerStart()
}
