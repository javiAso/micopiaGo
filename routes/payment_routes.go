package routes

import (
	"micopia/controllers"

	"github.com/gorilla/mux"
)

func SetPaymentRoutes(router *mux.Router) {
	subRoute := router.PathPrefix("/PaymentCRUD").Subrouter()
	subRoute.HandleFunc("/getPayments", controllers.AllPaymentsHandler).Methods("GET")
	subRoute.HandleFunc("/getCustomerPayments", controllers.GetCustomerPaymentsHandler).Methods("GET")
	subRoute.HandleFunc("/getPayment", controllers.GetPaymentHandler).Methods("GET")
	subRoute.HandleFunc("/createPayment", controllers.CreatePaymentHandler).Methods("PUT")
	subRoute.HandleFunc("/updatePayment", controllers.UpdatePaymentHandler).Methods("POST")
	subRoute.HandleFunc("/deletePayment", controllers.DeletePaymentHandler).Methods("DELETE")
}
