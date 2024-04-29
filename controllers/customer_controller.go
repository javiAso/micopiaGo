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
// @Router			/UserCRUD/getUsers [get]
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

//Private methods:
