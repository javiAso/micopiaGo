package routes

import (
	"micopia/controllers"

	"github.com/gorilla/mux"
)

func SetUserRoutes(router *mux.Router) {
	subRoute := router.PathPrefix("/UserCRUD").Subrouter()
	subRoute.HandleFunc("/getUsers", controllers.AllCustomersHandler).Methods("GET")
}
