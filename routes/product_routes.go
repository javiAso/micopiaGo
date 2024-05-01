package routes

import (
	"micopia/controllers"

	"github.com/gorilla/mux"
)

func SetProductRoutes(router *mux.Router) {
	subRoute := router.PathPrefix("/ProductCRUD").Subrouter()
	subRoute.HandleFunc("/getProducts", controllers.AllProductsHandler).Methods("GET")
	subRoute.HandleFunc("/getProduct", controllers.GetProductHandler).Methods("GET")
	subRoute.HandleFunc("/createProduct", controllers.CreateProductHandler).Methods("PUT")
	subRoute.HandleFunc("/updateProduct", controllers.UpdateProductHandler).Methods("POST")
	subRoute.HandleFunc("/deleteProduct", controllers.DeleteProductHandler).Methods("DELETE")
}
