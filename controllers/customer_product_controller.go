package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"micopia/commons"
	"micopia/models"
	"micopia/utils"
	"net/http"
	"strconv"
)

// CreateTags		godoc
// @Summary: 		Get All Customers Products
// @Description  	Get All Customers Products (cart/wishlist) from the database
// @Produce 		application/json
// @Tags			CustomerProduct
// @Success			200 {object} models.CustomersProducts
// @Router			/CustomerProductCRUD/getCustomersProducts [get]
func AllCustomersProductsHandler(w http.ResponseWriter, r *http.Request) {

	db := commons.GetConnection()

	var customersProducts []models.CustomerProduct
	rows, err := db.Query("SELECT c.customer_id, c.product_id, c.quantity, c.total_price FROM customer_product c;")
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var cP models.CustomerProduct
		if err := rows.Scan(&cP.Customer_id, &cP.Product_id, &cP.Quantity, &cP.TotalPrice); err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		customersProducts = append(customersProducts, cP)
	}

	if err := rows.Err(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var customerProductList models.CustomersProducts
	customerProductList.CustomerProductList = customersProducts

	// Convertir el resultado en formato JSON y enviar la respuesta
	utils.JSONResponse(w, http.StatusOK, customerProductList)

}

// CreateTags		godoc
// @Summary: 		Get CustomerProducts
// @Description  	Get CustomerProducts (cart/wishlist) from the database by customer id
// @Param			customerId query string true "The Customer identifier"
// @Produce 		application/json
// @Tags			CustomerProduct
// @Success			200 {object} []models.CustomerProduct
// @Router			/CustomerProductCRUD/getCustomerProduct [get]
func GetCustomerProductsHandler(w http.ResponseWriter, r *http.Request) {

	// Get the  "customerId"
	id := r.URL.Query().Get("customerId")
	if id == "" {
		utils.JSONError(w, http.StatusBadRequest, "missing customerId query parameter")
		return
	}

	// Get connection with the DB
	db := commons.GetConnection()
	defer db.Close()

	// Query to get customer product by id
	rows, err := db.Query("SELECT c.customer_id, c.product_id, c.quantity, c.total_price FROM customer_product c WHERE c.customer_id = ?;", id)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to get cart (DB QUERY ERROR)")
		return
	}
	defer rows.Close()

	//We scan the Data
	var c models.CustomerProduct
	var cP []models.CustomerProduct

	for rows.Next() {
		if err := rows.Scan(&c.Customer_id, &c.Product_id, &c.Quantity, &c.TotalPrice); err != nil {
			utils.JSONError(w, http.StatusInternalServerError, "failed to get cart ( ROWS SCAN ERROR)")
			return
		}
		cP = append(cP, c)
	}

	if err := rows.Err(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to get cart (ROWS ERROR)")
		return
	}

	//TODO: We are sending a null body here if there are no data

	// response as JSON
	utils.JSONResponse(w, http.StatusOK, cP)

}

// CreateTags		godoc
// @Summary: 		Create CustomerProduct
// @Description  	Create CustomerProduct in the database
// @Param			CreateCustomerProductRequest body models.CustomerProduct true "The CustomerProduct to create"
// @Produce 		application/json
// @Tags			CustomerProduct
// @Success      	201 {object} models.CustomerProduct
// @Router			/CustomerProductCRUD/createCustomerProduct [put]
func CreateCustomerProductHandler(w http.ResponseWriter, r *http.Request) {
	// Decodificar el cuerpo de la petición en una estructura CreateCustomerProductRequest
	var cP models.CustomerProduct
	if err := json.NewDecoder(r.Body).Decode(&cP); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed in customerProduct Request")
		return
	}

	// Obtener una conexión a la base de datos
	db := commons.GetConnection()
	defer db.Close()

	// Iniciar transacción en la base de datos
	tx, err := db.Begin()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to create customerProduct (Internal server error DB.BEGIN)")
		return
	}

	code, err := createCustomerProduct(&cP, tx)
	if err != nil {
		utils.JSONError(w, code, err.Error())
		return
	}

	// Realizar commit de la transacción en la base de datos
	if err := tx.Commit(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to create customerProduct (COMMIT FAILED)")
		return
	}

	// Devolver el nuevo CustomerProduct como JSON
	utils.JSONResponse(w, http.StatusCreated, cP)
}

// CreateTags		godoc
// @Summary: 		Update Customer Product
// @Description  	Update Customer (cart/wishlist) in the database
// @Param			Customer body models.CustomerProduct true "The Customer to update"
// @Produce 		application/json
// @Tags			CustomerProduct
// @Success      	200 {object} models.CustomerProduct
// @Router			/CustomerProductCRUD/updateCustomerProduct [post]
func UpdateCustomerProductHandler(w http.ResponseWriter, r *http.Request) {
	// Decodificar el cuerpo de la petición en una estructura Customer
	var c models.CustomerProduct
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed in customer product request")
		return
	}

	// Obtener una conexión a la base de datos
	db := commons.GetConnection()
	defer db.Close()

	// Iniciar transacción en la base de datos
	tx, err := db.Begin()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to update customer product (Internal server error DD.BEGIN)")
		return
	}

	code, err := updateCustomerProduct(&c, tx)
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
// @Summary: 		Delete CustomerProduct
// @Description  	Delete CustomerProduct in the database
// @Param			customerId query string true "The Customer identifier"
// @Param			productId query string true "The Product identifier"
// @Produce 		application/json
// @Tags			CustomerProduct
// @Success      	200 {string} string "Customer product deleted successfully"
// @Router 			/CustomerProductCRUD/deleteCustomerProduct [delete]
func DeleteCustomerProductHandler(w http.ResponseWriter, r *http.Request) {
	// Get a connection to the database
	db := commons.GetConnection()
	defer db.Close()

	// Extract the id from the URL segment

	cId, err := strconv.Atoi(r.URL.Query().Get("customerId"))
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	pId, err := strconv.Atoi(r.URL.Query().Get("productId"))
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
	code, err := deleteCustomerProduct(cId, pId, tx)
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
	utils.JSONResponse(w, http.StatusOK, "Customer product deleted successfully")
}

//Private methods:

func createCustomerProduct(cP *models.CustomerProduct, tx *sql.Tx) (int, error) {
	//Check the required fields
	if cP.Customer_id == 0 {
		return http.StatusBadRequest, errors.New("customer_id can not be empty")
	}
	if cP.Product_id == 0 {
		return http.StatusBadRequest, errors.New("product_id can not be empty")
	}

	// Se prepara la sentencia SQL para insertar la hora
	stmt, err := tx.Prepare("INSERT INTO customer_product(customer_id, product_id, quantity, total_price) VALUES(?, ?, ?, ?)")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer stmt.Close()

	// We execute the insert
	result, err := stmt.Exec(cP.Customer_id, cP.Product_id, cP.Quantity, cP.TotalPrice)
	if err != nil {
		tx.Rollback()
		return http.StatusNotAcceptable, err
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
		return http.StatusNotModified, errors.New("product not created")
	}

	return http.StatusCreated, nil
}

func updateCustomerProduct(c *models.CustomerProduct, tx *sql.Tx) (int, error) {
	//Check the required fields

	if c.Customer_id == 0 {
		return http.StatusBadRequest, errors.New("customer_id can not be empty")
	}

	//Check the required fields
	if c.Product_id == 0 {
		return http.StatusBadRequest, errors.New("product_id can not be empty")
	}

	// Se prepara la sentencia SQL para insertar el Customer
	stmt, err := tx.Prepare("UPDATE customer_product SET quantity = ?, total_price = ? WHERE customer_id = ? AND product_id = ?")
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	defer stmt.Close()

	// Se ejecuta la sentencia SQL para updatear el Customer
	result, err := stmt.Exec(c.Quantity, c.TotalPrice, c.Customer_id, c.Product_id)
	if err != nil {
		tx.Rollback()
		return http.StatusNotAcceptable, err
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
		return http.StatusNotModified, errors.New("customer_product not updated (NO CHANGES)")
	}

	return http.StatusOK, nil
}

func deleteCustomerProduct(cId int, pId int, tx *sql.Tx) (int, error) {

	// Prepare the SQL statement to delete the customer
	stmt, err := tx.Prepare("DELETE FROM customer_product WHERE customer_id = ? AND product_id = ?")
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to delete customer (PREPARE QUERY)")
	}
	defer stmt.Close()

	// Execute the SQL statement to delete the customer
	result, err := stmt.Exec(cId, pId)
	if err != nil {
		tx.Rollback()
		return http.StatusNotAcceptable, errors.New("failed to delete customer product (EXECUTE QUERY)")
	}

	// Get the number of rows affected by the SQL statement
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to delete customer product (ROWS AFFECTED)")
	}

	// If no rows were deleted, rollback the transaction
	if rowsAffected == 0 {
		tx.Rollback()
		return http.StatusNotModified, errors.New("failed to delete customer product (NO ROWS AFFECTED)")
	}
	return http.StatusOK, nil
}
