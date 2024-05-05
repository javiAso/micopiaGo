package routes

import (
	"micopia/controllers"

	"github.com/gorilla/mux"
)

func SetShipmentRoutes(router *mux.Router) {
	subRoute := router.PathPrefix("/ShipmentCRUD").Subrouter()
	subRoute.HandleFunc("/getShipments", controllers.AllShipmentsHandler).Methods("GET")
	subRoute.HandleFunc("/getCustomerShipments", controllers.GetCustomerShipmentsHandler).Methods("GET")
	subRoute.HandleFunc("/getShipment", controllers.GetShipmentHandler).Methods("GET")
	subRoute.HandleFunc("/createShipment", controllers.CreateShipmentHandler).Methods("PUT")
	subRoute.HandleFunc("/updateShipment", controllers.UpdateShipmentHandler).Methods("POST")
	subRoute.HandleFunc("/deleteShipment", controllers.DeleteShipmentHandler).Methods("DELETE")
}
