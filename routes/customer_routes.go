package routes

import (
	"micopia/controllers"

	"github.com/gorilla/mux"
)

func SetCustomerRoutes(router *mux.Router) {
	subRoute := router.PathPrefix("/CustomerCRUD").Subrouter()
	subRoute.HandleFunc("/getCustomers", controllers.AllCustomersHandler).Methods("GET")
	subRoute.HandleFunc("/getCustomer", controllers.GetCustomerHandler).Methods("GET")
	subRoute.HandleFunc("/createCustomer", controllers.CreateCustomerHandler).Methods("PUT")
	subRoute.HandleFunc("/updateCustomer", controllers.UpdateCustomerHandler).Methods("POST")

}
