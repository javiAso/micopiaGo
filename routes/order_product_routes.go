package routes

import (
	"micopia/controllers"

	"github.com/gorilla/mux"
)

func SetOrderProductRoutes(router *mux.Router) {
	subRoute := router.PathPrefix("/OrderProductCRUD").Subrouter()
	subRoute.HandleFunc("/getOrdersProducts", controllers.AllOrdersProductsHandler).Methods("GET")
	subRoute.HandleFunc("/getOrderProducts", controllers.GetOrderProductsHandler).Methods("GET")
	subRoute.HandleFunc("/getProductOrders", controllers.GetProductOrdersHandler).Methods("GET")
	subRoute.HandleFunc("/createOrderProduct", controllers.CreateOrderProductHandler).Methods("PUT")
	subRoute.HandleFunc("/updateOrderProduct", controllers.UpdateOrderProductHandler).Methods("POST")
	subRoute.HandleFunc("/deleteOrderProduct", controllers.DeleteOrderProductHandler).Methods("DELETE")
}
