package routes

import (
	"micopia/controllers"

	"github.com/gorilla/mux"
)

func SetOrderRoutes(router *mux.Router) {
	subRoute := router.PathPrefix("/OrderCRUD").Subrouter()
	subRoute.HandleFunc("/getOrders", controllers.AllOrdersHandler).Methods("GET")
	subRoute.HandleFunc("/getCustomerOrders", controllers.GetCustomerOrdersHandler).Methods("GET")
	subRoute.HandleFunc("/getPaymentOrders", controllers.GetPaymentOrdersHandler).Methods("GET")
	subRoute.HandleFunc("/getShipmentOrders", controllers.GetShipmentOrdersHandler).Methods("GET")
	subRoute.HandleFunc("/getOrder", controllers.GetOrderHandler).Methods("GET")
	subRoute.HandleFunc("/createOrder", controllers.CreateOrderHandler).Methods("PUT")
	subRoute.HandleFunc("/updateOrder", controllers.UpdateOrderHandler).Methods("POST")
	subRoute.HandleFunc("/deleteOrder", controllers.DeleteOrderHandler).Methods("DELETE")
}
