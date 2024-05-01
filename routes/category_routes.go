package routes

import (
	"micopia/controllers"

	"github.com/gorilla/mux"
)

func SetCategoryRoutes(router *mux.Router) {
	subRoute := router.PathPrefix("/categoryCRUD").Subrouter()
	subRoute.HandleFunc("/getCategories", controllers.AllCategorysHandler).Methods("GET")
	subRoute.HandleFunc("/getCategory", controllers.GetCategoryHandler).Methods("GET")
	/*  	subRoute.HandleFunc("/createcategory", controllers.CreatecategoryHandler).Methods("PUT")
	subRoute.HandleFunc("/updatecategory", controllers.UpdatecategoryHandler).Methods("POST")
	subRoute.HandleFunc("/deletecategory", controllers.DeletecategoryHandler).Methods("DELETE") */
}
