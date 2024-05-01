package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"micopia/commons"
	"micopia/models"
	"micopia/utils"
	"net/http"
	"net/mail"
	"regexp"
	"strconv"

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

// CreateTags		godoc
// @Summary: 		Create Customer
// @Description  	Create Customer in the database
// @Param			CreateCustomerRequest body models.CreateCustomerRequest true "The Customer to create"
// @Produce 		application/json
// @Tags			Customer
// @Success      	201 {object} models.Customer
// @Router			/CustomerCRUD/createCustomer [put]
func CreateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	// Decodificar el cuerpo de la petición en una estructura CreateCustomerRequest
	var createReq models.CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Failed in Customer Request")
		return
	}

	// Obtener una conexión a la base de datos
	db := commons.GetConnection()
	defer db.Close()

	// Iniciar transacción en la base de datos
	tx, err := db.Begin()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Failed to create Customer (Internal server error DB.BEGIN)")
		return
	}

	customer := mapCreateCustomerReqToCustomer(createReq)

	code, err := createCustomer(&customer, tx)
	if err != nil {
		utils.JSONError(w, code, err.Error())
		return
	}

	// Realizar commit de la transacción en la base de datos
	if err := tx.Commit(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Failed to create Customer (COMMIT FAILED)")
		return
	}

	// Devolver el nuevo Customer como JSON
	utils.JSONResponse(w, http.StatusCreated, customer)
}

// CreateTags		godoc
// @Summary: 		Update Customer
// @Description  	Update Customer in the database
// @Param			UpdateCustomerRequest body models.Customer true "The Customer to update"
// @Produce 		application/json
// @Tags			Customer
// @Success      	200 {object} models.Customer
// @Router			/CustomerCRUD/updateCustomer [post]
func UpdateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	// Decodificar el cuerpo de la petición en una estructura Customer
	var c models.Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed in customer request")
		return
	}

	// Obtener una conexión a la base de datos
	db := commons.GetConnection()
	defer db.Close()

	// Iniciar transacción en la base de datos
	tx, err := db.Begin()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to update customer (Internal server error DD.BEGIN)")
		return
	}

	code, err := updateCustomer(&c, tx)
	if err != nil {
		utils.JSONError(w, code, err.Error())
		return
	}

	// Realizar commit de la transacción en la base de datos
	if err := tx.Commit(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to update customer (COMMIT FAILED)")
		return
	}

	// Devolver el nuevo Customer como JSON
	utils.JSONResponse(w, http.StatusOK, c)
}

// CreateTags		godoc
// @Summary: 		Delete Customer
// @Description  	Delete Customer in the database
// @Param			customerId query string true "The Customer identifier"
// @Produce 		application/json
// @Tags			Customer
// @Success      	200 {string} string "Customer deleted successfully"
// @Router 			/CustomerCRUD/deleteCustomer [delete]
func DeleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
	// Get a connection to the database
	db := commons.GetConnection()
	defer db.Close()

	// Extract the id from the URL segment

	id, err := strconv.Atoi(r.URL.Query().Get("customerId"))
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	// Start a transaction in the database
	tx, err := db.Begin()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Failed to delete Customer (CONNECTING DB)")
		return
	}
	code, err := deleteCustomer(id, tx)
	if err != nil {
		utils.JSONError(w, code, err.Error())
		return
	}
	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Failed to commit the delete transaction")
		return
	}
	// Return a successful response
	utils.JSONResponse(w, http.StatusOK, "Customer deleted successfully")
}

//Private methods:

func createCustomer(c *models.Customer, tx *sql.Tx) (int, error) {
	//Check the required fields
	if c.Email == "" {
		return http.StatusBadRequest, errors.New("email can not be empty")
	}
	if c.First_name == "" {
		return http.StatusBadRequest, errors.New("first name can not be empty")
	}

	if c.Last_name == "" {
		return http.StatusBadRequest, errors.New("last name can not be empty")
	}

	if c.Address == "" {
		return http.StatusBadRequest, errors.New("address can not be empty")
	}

	if c.Phone_Number == "" {
		return http.StatusBadRequest, errors.New("phone number can not be empty")
	}

	_, err := mail.ParseAddress(c.Email)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if !regexp.MustCompile(`^\+?[1-9]\d{1,14}$`).MatchString(c.Phone_Number) {
		return http.StatusBadRequest, errors.New("phone number is not valid")
	}

	// We prepare the query to create a user
	stmt, err := tx.Prepare("INSERT INTO Customer(address, email, first_name, last_name, phone_number) VALUES(?, ?, ?, ?,?)")
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	defer stmt.Close()

	// Se ejecuta la sentencia SQL para insertar el user
	result, err := stmt.Exec(c.Address, c.Email, c.First_name, c.Last_name, c.Phone_Number)
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}

	// Se obtiene la cantidad de filas afectadas por la sentencia SQL
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	// Si no se ha insertado ninguna fila, se hace rollback de la transacción
	if rowsAffected == 0 {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("customer not created")
	}

	// Get the ID of the newly created user
	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}

	// Set the ID of the user struct to the new ID
	c.Customer_id = uint64(id)

	return http.StatusCreated, nil
}

func updateCustomer(c *models.Customer, tx *sql.Tx) (int, error) {
	//Check the required fields

	if c.Customer_id == 0 {
		return http.StatusBadRequest, errors.New("customer_id can not be empty")
	}

	//Check the required fields
	if c.Email == "" {
		return http.StatusBadRequest, errors.New("email can not be empty")
	}
	if c.First_name == "" {
		return http.StatusBadRequest, errors.New("first name can not be empty")
	}

	if c.Last_name == "" {
		return http.StatusBadRequest, errors.New("last name can not be empty")
	}

	if c.Address == "" {
		return http.StatusBadRequest, errors.New("address can not be empty")
	}

	if c.Phone_Number == "" {
		return http.StatusBadRequest, errors.New("phone number can not be empty")
	}

	_, err := mail.ParseAddress(c.Email)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if !regexp.MustCompile(`^\+?[1-9]\d{1,14}$`).MatchString(c.Phone_Number) {
		return http.StatusBadRequest, errors.New("phone number is not valid")
	}

	// Se prepara la sentencia SQL para insertar el Customer
	stmt, err := tx.Prepare("UPDATE Customer SET email = ?, first_name = ?, last_name = ?, phone_number = ?, address = ? WHERE customer_id = ?")
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	defer stmt.Close()

	// Se ejecuta la sentencia SQL para updatear el Customer
	result, err := stmt.Exec(c.Email, c.First_name, c.Last_name, c.Phone_Number, c.Address, c.Customer_id)
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}

	// Se obtiene la cantidad de filas afectadas por la sentencia SQL
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	// Si no se ha updateado ninguna fila, se hace rollback de la transacción
	if rowsAffected == 0 {
		tx.Rollback()
		return http.StatusNotModified, errors.New("customer not updated (NO CHANGES)")
	}

	return http.StatusOK, nil
}

func deleteCustomer(id int, tx *sql.Tx) (int, error) {

	// Prepare the SQL statement to delete the customer
	stmt, err := tx.Prepare("DELETE FROM Customer WHERE customer_id = ?")
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to delete customer (PREPARE QUERY)")
	}
	defer stmt.Close()

	// Execute the SQL statement to delete the customer
	result, err := stmt.Exec(id)
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to delete customer (EXECUTE QUERY)")
	}

	// Get the number of rows affected by the SQL statement
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to delete customer (ROWS AFFECTED)")
	}

	// If no rows were deleted, rollback the transaction
	if rowsAffected == 0 {
		tx.Rollback()
		return http.StatusNotModified, errors.New("failed to delete customer (NO ROWS AFFECTED)")
	}
	return http.StatusOK, nil
}

func mapCreateCustomerReqToCustomer(c models.CreateCustomerRequest) models.Customer {
	return models.Customer{
		First_name:   c.First_name,
		Last_name:    c.Last_name,
		Email:        c.Email,
		Address:      c.Address,
		Phone_Number: c.Phone_Number,
	}
}
