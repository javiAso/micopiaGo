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

	_ "github.com/go-sql-driver/mysql"
)

// CreateTags		godoc
// @Summary: 		Get All Payments
// @Description  	Get All Payments from the database
// @Produce 		application/json
// @Tags			Payment
// @Success			200 {object} []models.Payment
// @Router			/PaymentCRUD/getPayments [get]
func AllPaymentsHandler(w http.ResponseWriter, r *http.Request) {

	db := commons.GetConnection()

	var payments []models.Payment
	rows, err := db.Query("SELECT p.payment_id, p.payment_date, p.payment_method, p.amount, p.customer_id FROM payment p;")
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var payment models.Payment
		if err := rows.Scan(&payment.Payment_id, &payment.Payment_date, &payment.Payment_method, &payment.Amount, &payment.Customer_id); err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// We return Json response:
	utils.JSONResponse(w, http.StatusOK, payments)

}

// CreateTags		godoc
// @Summary: 		Get All Payments by Customer Id
// @Description  	Get All Payments by Customer Id from the database
// @Param			customerId query string true "The Customer identifier"
// @Produce 		application/json
// @Tags			Payment
// @Success			200 {object} []models.Payment
// @Router			/PaymentCRUD/getCustomerPayments [get]
func GetCustomerPaymentsHandler(w http.ResponseWriter, r *http.Request) {

	// We get "customerId" from URL
	id := r.URL.Query().Get("customerId")
	if id == "" {
		utils.JSONError(w, http.StatusBadRequest, "missing customerId query parameter")
		return
	}

	db := commons.GetConnection()

	var payments []models.Payment
	rows, err := db.Query("SELECT p.payment_id, p.payment_date, p.payment_method, p.amount, p.customer_id FROM payment p WHERE p.customer_id = ?;", id)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var payment models.Payment
		if err := rows.Scan(&payment.Payment_id, &payment.Payment_date, &payment.Payment_method, &payment.Amount, &payment.Customer_id); err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// We return Json response:
	utils.JSONResponse(w, http.StatusOK, payments)

}

// CreateTags		godoc
// @Summary: 		Get Payment by Payment Id
// @Description  	Get Payment by Payment Id from the database
// @Param			paymentId query string true "The Payment identifier"
// @Produce 		application/json
// @Tags			Payment
// @Success			200 {object} models.Payment
// @Router			/PaymentCRUD/getPayment [get]
func GetPaymentHandler(w http.ResponseWriter, r *http.Request) {

	// We get "paymentId" from URL
	id := r.URL.Query().Get("paymentId")
	if id == "" {
		utils.JSONError(w, http.StatusBadRequest, "missing paymentId query parameter")
		return
	}

	db := commons.GetConnection()

	rows, err := db.Query("SELECT p.payment_id, p.payment_date, p.payment_method, p.amount, p.customer_id FROM payment p WHERE p.payment_id = ?;", id)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	// We  set the data
	var payment models.Payment
	for rows.Next() {
		if err := rows.Scan(&payment.Payment_id, &payment.Payment_date, &payment.Payment_method, &payment.Amount, &payment.Customer_id); err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if err := rows.Err(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// We return Json response:
	utils.JSONResponse(w, http.StatusOK, payment)

}

// CreateTags		godoc
// @Summary: 		Create Payment
// @Description  	Create Payment in the database
// @Param			CreatePaymentRequest body models.Payment true "The Payment to create, payment_id is not relevant"
// @Produce 		application/json
// @Tags			Payment
// @Success      	201 {object} models.Payment
// @Router			/PaymentCRUD/createPayment [put]
func CreatePaymentHandler(w http.ResponseWriter, r *http.Request) {
	// Decodificar el cuerpo de la petición en una estructura CreatePaymentRequest
	var createReq models.Payment
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed in payment Request")
		return
	}

	// Obtener una conexión a la base de datos
	db := commons.GetConnection()
	defer db.Close()

	// Iniciar transacción en la base de datos
	tx, err := db.Begin()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to create payment (Internal server error DB.BEGIN)")
		return
	}

	code, err := createPayment(&createReq, tx)
	if err != nil {
		utils.JSONError(w, code, err.Error())
		return
	}

	// Realizar commit de la transacción en la base de datos
	if err := tx.Commit(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to create payment (COMMIT FAILED)")
		return
	}

	// Devolver el nuevo Payment como JSON
	utils.JSONResponse(w, http.StatusCreated, createReq)
}

// CreateTags		godoc
// @Summary: 		Update Payment
// @Description  	Update Payment in the database
// @Param			UpdatePaymentRequest body models.Payment true "The Payment to update"
// @Produce 		application/json
// @Tags			Payment
// @Success      	200 {object} models.Payment
// @Router			/PaymentCRUD/updatePayment [post]
func UpdatePaymentHandler(w http.ResponseWriter, r *http.Request) {
	// Decodificar el cuerpo de la petición en una estructura Payment
	var p models.Payment
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed in payment request")
		return
	}

	// Obtener una conexión a la base de datos
	db := commons.GetConnection()
	defer db.Close()

	// Iniciar transacción en la base de datos
	tx, err := db.Begin()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to update payment (Internal server error DD.BEGIN)")
		return
	}

	code, err := updatePayment(&p, tx)
	if err != nil {
		utils.JSONError(w, code, err.Error())
		return
	}

	// Realizar commit de la transacción en la base de datos
	if err := tx.Commit(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to update payment (COMMIT FAILED)")
		return
	}

	// Devolver el nuevo Payment como JSON
	utils.JSONResponse(w, http.StatusOK, p)
}

// CreateTags		godoc
// @Summary: 		Delete Payment
// @Description  	Delete Payment in the database
// @Param			paymentId query string true "The Payment identifier"
// @Produce 		application/json
// @Tags			Payment
// @Success      	200 {string} string "Payment deleted successfully"
// @Router 			/PaymentCRUD/deletePayment [delete]
func DeletePaymentHandler(w http.ResponseWriter, r *http.Request) {
	// Get a connection to the database
	db := commons.GetConnection()
	defer db.Close()

	// Extract the id from the URL segment

	id, err := strconv.Atoi(r.URL.Query().Get("paymentId"))
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	// Start a transaction in the database
	tx, err := db.Begin()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to delete payment (CONNECTING DB)")
		return
	}
	code, err := deletePayment(id, tx)
	if err != nil {
		utils.JSONError(w, code, err.Error())
		return
	}
	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to commit the delete transaction")
		return
	}
	// Return a successful response
	utils.JSONResponse(w, http.StatusOK, "payment deleted successfully")
}

// private methods

func createPayment(p *models.Payment, tx *sql.Tx) (int, error) {
	//Check the required fields

	if p.Customer_id == 0 {
		return http.StatusBadRequest, errors.New("customer_id can not be empty")
	}

	if p.Payment_method == "" {
		return http.StatusBadRequest, errors.New("payment_method can not be empty")

	}

	if p.Payment_date == "" { //TODO: use time library for check if we have a valid DATE
		return http.StatusBadRequest, errors.New("payment_date can not be empty")
	}

	// We prepare the query to create a payment
	stmt, err := tx.Prepare("INSERT INTO Payment(payment_date, payment_method, amount, customer_id) VALUES(?,?,?,?)")
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	defer stmt.Close()

	// We execute SQL to insert the payment
	result, err := stmt.Exec(p.Payment_date, p.Payment_method, p.Amount, p.Customer_id)
	if err != nil {
		tx.Rollback()
		return http.StatusNotAcceptable, err
	}

	// Se obtiene la cantidad de filas afectadas por la sentencia SQL
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return http.StatusNotModified, err
	}
	// Si no se ha insertado ninguna fila, se hace rollback de la transacción
	if rowsAffected == 0 {
		tx.Rollback()
		return http.StatusNotModified, errors.New("payment not created")
	}

	// Get the ID of the newly created user
	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return http.StatusNotModified, err
	}

	// Set the ID of the user struct to the new ID
	p.Payment_id = uint64(id)

	return http.StatusCreated, nil
}

func updatePayment(p *models.Payment, tx *sql.Tx) (int, error) {
	//Check the required fields

	if p.Payment_id == 0 {
		return http.StatusBadRequest, errors.New("payment_id can not be empty")
	}

	if p.Payment_method == "" {
		return http.StatusBadRequest, errors.New("payment_method can not be empty")

	}

	if p.Payment_date == "" { //TODO: use time library for check if we have a valid DATE
		return http.StatusBadRequest, errors.New("payment_date can not be empty")
	}

	// Se prepara la sentencia SQL para updatear el Payment
	stmt, err := tx.Prepare("UPDATE payment SET payment_method = ?, payment_date = ?, amount = ? WHERE payment_id = ?")
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	defer stmt.Close()

	// Se ejecuta la sentencia SQL para updatear el Payment
	result, err := stmt.Exec(p.Payment_method, p.Payment_date, p.Amount, p.Payment_id)
	if err != nil {
		tx.Rollback()
		return http.StatusNotAcceptable, err
	}

	// Se obtiene la cantidad de filas afectadas por la sentencia SQL
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return http.StatusNotModified, err
	}
	// Si no se ha updateado ninguna fila, se hace rollback de la transacción
	if rowsAffected == 0 {
		tx.Rollback()
		return http.StatusNotModified, errors.New("payment not updated (NO CHANGES)")
	}

	return http.StatusOK, nil
}

func deletePayment(id int, tx *sql.Tx) (int, error) {

	// Prepare the SQL statement to delete the payment
	stmt, err := tx.Prepare("DELETE FROM Payment WHERE payment_id = ?")
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to delete payment (PREPARE QUERY)")
	}
	defer stmt.Close()

	// Execute the SQL statement to delete the payment
	result, err := stmt.Exec(id)
	if err != nil {
		tx.Rollback()
		return http.StatusNotAcceptable, errors.New("failed to delete payment (EXECUTE QUERY)")
	}

	// Get the number of rows affected by the SQL statement
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to delete payment (ROWS AFFECTED)")
	}

	// If no rows were deleted, rollback the transaction
	if rowsAffected == 0 {
		tx.Rollback()
		return http.StatusNotModified, errors.New("failed to delete payment (NO ROWS AFFECTED)")
	}
	return http.StatusOK, nil
}
