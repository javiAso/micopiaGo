package controllers

import (
	"micopia/commons"
	"micopia/models"
	"micopia/utils"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// CreateTags		godoc
// @Summary: 		Get All Customers
// @Description  	Get All Customers from the database
// @Produce 		application/json
// @Tags			Customer
// @Success			200 {object} models.Customers
// @Router			/CustomerCRUD/getCustomers [get]
func AllCustomersHandler(w http.ResponseWriter, r *http.Request) {

	db := commons.GetConnection()

	var customers []models.Customer
	rows, err := db.Query("SELECT c.customer_id, c.first_name, c.last_name, c.email, c.address, c.phone_number FROM customer c;")
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var customer models.Customer
		if err := rows.Scan(&customer.Customer_id, &customer.First_name, &customer.Last_name, &customer.Email, &customer.Address, &customer.Phone_Number); err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		customers = append(customers, customer)
	}

	if err := rows.Err(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var customerList models.Customers
	customerList.CustomerList = customers

	// Convertir el resultado en formato JSON y enviar la respuesta
	utils.JSONResponse(w, http.StatusOK, customerList)

}

// CreateTags		godoc
// @Summary: 		Get Customer
// @Description  	Get Customer from the database by id
// @Param			customerId query string true "The Customer identifier"
// @Produce 		application/json
// @Tags			Customer
// @Success			200 {object} models.Customer
// @Router			/CustomerCRUD/getCustomer [get]
func GetCustomerHandler(w http.ResponseWriter, r *http.Request) {

	// Obtener el ID del proyecto del parámetro de consulta "customerId"
	id := r.URL.Query().Get("customerId")
	if id == "" {
		utils.JSONError(w, http.StatusBadRequest, "Missing customerId query parameter")
		return
	}

	// Obtener una conexión a la base de datos
	db := commons.GetConnection()
	defer db.Close()

	// Realizar la consulta SQL para obtener las horas del proyecto
	rows, err := db.Query("SELECT c.customer_id, c.first_name, c.last_name, c.email, c.address, c.phone_number FROM customer c WHERE c.customer_id = ?;", id)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to get customer (DB QUERY ERROR)")
		return
	}
	defer rows.Close()

	//We scan the Data
	var c models.Customer
	for rows.Next() {
		if err := rows.Scan(&c.Customer_id, &c.First_name, &c.Last_name, &c.Email, &c.Address, &c.Phone_Number); err != nil {
			utils.JSONError(w, http.StatusInternalServerError, "failed to get customer ( ROWS SCAN ERROR)")
			return
		}
	}

	if err := rows.Err(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to get customer (ROWS ERROR)")
		return
	}

	// Devolver la nueva organización como JSON
	utils.JSONResponse(w, http.StatusCreated, c)

}

//Private methods:
