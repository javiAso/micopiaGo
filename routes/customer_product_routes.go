package routes

import (
	"micopia/controllers"

	"github.com/gorilla/mux"
)

func SetCustomerProductRoutes(router *mux.Router) {
	subRoute := router.PathPrefix("/CustomerProductCRUD").Subrouter()
	subRoute.HandleFunc("/getCustomersProducts", controllers.AllCustomersProductsHandler).Methods("GET")
	subRoute.HandleFunc("/getCustomerProduct", controllers.GetCustomerProductsHandler).Methods("GET")
	subRoute.HandleFunc("/createCustomerProduct", controllers.CreateCustomerProductHandler).Methods("PUT")
	subRoute.HandleFunc("/updateCustomerProduct", controllers.UpdateCustomerProductHandler).Methods("POST")
	subRoute.HandleFunc("/deleteCustomerProduct", controllers.DeleteCustomerProductHandler).Methods("DELETE")
}
