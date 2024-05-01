package routes

import (
	"micopia/controllers"

	"github.com/gorilla/mux"
)

func SetCategoryRoutes(router *mux.Router) {
	subRoute := router.PathPrefix("/CategoryCRUD").Subrouter()
	subRoute.HandleFunc("/getCategories", controllers.AllCategorysHandler).Methods("GET")
	subRoute.HandleFunc("/getCategory", controllers.GetCategoryHandler).Methods("GET")
	subRoute.HandleFunc("/createCategory", controllers.CreateCategoryHandler).Methods("PUT")
	subRoute.HandleFunc("/updateCategory", controllers.UpdateCategoryHandler).Methods("POST")
	subRoute.HandleFunc("/deleteCategory", controllers.DeleteCategoryHandler).Methods("DELETE")
}
