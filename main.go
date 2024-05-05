package main

import (
	"log"
	"micopia/routes"
	"net/http"
	"os"

	_ "micopia/docs"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func initControllers(router *mux.Router) {
	routes.SetCustomerRoutes(router)
	routes.SetProductRoutes(router)
	routes.SetCategoryRoutes(router)
	routes.SetCustomerProductRoutes(router)
	routes.SetPaymentRoutes(router)
	routes.SetShipmentRoutes(router)
}

func initSwagger(r *mux.Router) {
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)
}

// @title Micopia Swagger Documentation
// @version 1.0
// @description Wellcome to the Micopia Web Server Swagger Documentation
// @host localhost:8080

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar archivo .env")
	}

	// Crear un enrutador multiplexor
	router := mux.NewRouter()
	c := cors.AllowAll() // Permite solicitudes CORS desde cualquier origen
	initSwagger(router)
	initControllers(router)
	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":8080", handler))
	server := http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: router,
	}

	log.Println("Servidor iniciado")

	log.Println(server.ListenAndServe())

}
