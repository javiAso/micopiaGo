package routes

import (
	"micopia/controllers"

	"github.com/gorilla/mux"
)

func SetProductRoutes(router *mux.Router) {
	subRoute := router.PathPrefix("/ProductCRUD").Subrouter()
	subRoute.HandleFunc("/getProducts", controllers.AllProductsHandler).Methods("GET")
	subRoute.HandleFunc("/createProduct", controllers.CreateProduct).Methods("POST")
}
