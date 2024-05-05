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
// @Summary: 		Get All Shipments
// @Description  	Get All Shipments from the database
// @Produce 		application/json
// @Tags			Shipment
// @Success			200 {object} []models.Shipment
// @Router			/ShipmentCRUD/getShipments [get]
func AllShipmentsHandler(w http.ResponseWriter, r *http.Request) {

	db := commons.GetConnection()

	var shipments []models.Shipment
	rows, err := db.Query("SELECT s.shipment_id, s.shipment_date, s.customer_id, s.address, s.city, s.State, s.country, s.zip_code FROM shipment s;")
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var shipment models.Shipment
		if err := rows.Scan(&shipment.Shipment_id, &shipment.Shipment_date, &shipment.Customer_id, &shipment.Address, &shipment.City, &shipment.State, &shipment.Country, &shipment.Zip_code); err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		shipments = append(shipments, shipment)
	}

	if err := rows.Err(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// We return Json response:
	utils.JSONResponse(w, http.StatusOK, shipments)

}

// CreateTags		godoc
// @Summary: 		Get Customer Shipments
// @Description  	Get All Customer Shipments from the database
// @Param			customerId query string true "The customer identifier"
// @Produce 		application/json
// @Tags			Shipment
// @Success			200 {object} []models.Shipment
// @Router			/ShipmentCRUD/getCustomerShipments [get]
func GetCustomerShipmentsHandler(w http.ResponseWriter, r *http.Request) {

	// We get "customerId" from URL
	id := r.URL.Query().Get("customerId")
	if id == "" {
		utils.JSONError(w, http.StatusBadRequest, "missing customerId query parameter")
		return
	}

	db := commons.GetConnection()

	var shipments []models.Shipment
	rows, err := db.Query("SELECT s.shipment_id, s.shipment_date, s.customer_id, s.address, s.city, s.State, s.country, s.zip_code FROM shipment s WHERE s.customer_id = ?;", id)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var shipment models.Shipment
		if err := rows.Scan(&shipment.Shipment_id, &shipment.Shipment_date, &shipment.Customer_id, &shipment.Address, &shipment.City, &shipment.State, &shipment.Country, &shipment.Zip_code); err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		shipments = append(shipments, shipment)
	}

	if err := rows.Err(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// We return Json response:
	utils.JSONResponse(w, http.StatusOK, shipments)

}

// CreateTags		godoc
// @Summary: 		Get Shipment
// @Description  	Get Shipment from the database by id
// @Param			shipmentId query string true "The shipment identifier"
// @Produce 		application/json
// @Tags			Shipment
// @Success			200 {object} models.Shipment
// @Router			/ShipmentCRUD/getShipment [get]
func GetShipmentHandler(w http.ResponseWriter, r *http.Request) {

	// We get "shipmentId" from URL
	id := r.URL.Query().Get("shipmentId")
	if id == "" {
		utils.JSONError(w, http.StatusBadRequest, "missing shipmentId query parameter")
		return
	}

	db := commons.GetConnection()

	rows, err := db.Query("SELECT s.shipment_id, s.shipment_date, s.customer_id, s.address, s.city, s.State, s.country, s.zip_code FROM shipment s WHERE s.shipment_id = ?;", id)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var shipment models.Shipment
	for rows.Next() {
		if err := rows.Scan(&shipment.Shipment_id, &shipment.Shipment_date, &shipment.Customer_id, &shipment.Address, &shipment.City, &shipment.State, &shipment.Country, &shipment.Zip_code); err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if err := rows.Err(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// We return Json response:
	utils.JSONResponse(w, http.StatusOK, shipment)

}

// CreateTags		godoc
// @Summary: 		Create Shipment
// @Description  	Create Shipment in the database
// @Param			CreateShipmentRequest body models.Shipment true "The Shipment to create, shipment_id is not relevant"
// @Produce 		application/json
// @Tags			Shipment
// @Success      	201 {object} models.Shipment
// @Router			/ShipmentCRUD/createShipment [put]
func CreateShipmentHandler(w http.ResponseWriter, r *http.Request) {
	// Decodificar el cuerpo de la petición en una estructura CreateShipmentRequest
	var createReq models.Shipment
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed in shipment Request")
		return
	}

	// Obtener una conexión a la base de datos
	db := commons.GetConnection()
	defer db.Close()

	// Iniciar transacción en la base de datos
	tx, err := db.Begin()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to create shipment (Internal server error DB.BEGIN)")
		return
	}

	code, err := createShipment(&createReq, tx)
	if err != nil {
		utils.JSONError(w, code, err.Error())
		return
	}

	// Realizar commit de la transacción en la base de datos
	if err := tx.Commit(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to create shipment (COMMIT FAILED)")
		return
	}

	// Devolver el nuevo Shipment como JSON
	utils.JSONResponse(w, http.StatusCreated, createReq)
}

// CreateTags		godoc
// @Summary: 		Update Shipment
// @Description  	Update Shipment in the database
// @Param			UpdateShipmentRequest body models.Shipment true "The Shipment to update"
// @Produce 		application/json
// @Tags			Shipment
// @Success      	200 {object} models.Shipment
// @Router			/ShipmentCRUD/updateShipment [post]
func UpdateShipmentHandler(w http.ResponseWriter, r *http.Request) {
	// Decodificar el cuerpo de la petición en una estructura Shipment
	var s models.Shipment
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed in shipment request")
		return
	}

	// Obtener una conexión a la base de datos
	db := commons.GetConnection()
	defer db.Close()

	// Iniciar transacción en la base de datos
	tx, err := db.Begin()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to update shipment (Internal server error DD.BEGIN)")
		return
	}

	code, err := updateShipment(&s, tx)
	if err != nil {
		utils.JSONError(w, code, err.Error())
		return
	}

	// Realizar commit de la transacción en la base de datos
	if err := tx.Commit(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to update shipment (COMMIT FAILED)")
		return
	}

	// Devolver el nuevo Shipment como JSON
	utils.JSONResponse(w, http.StatusOK, s)
}

// CreateTags		godoc
// @Summary: 		Delete Shipment
// @Description  	Delete Shipment in the database
// @Param			shipmentId query string true "The Shipment identifier"
// @Produce 		application/json
// @Tags			Shipment
// @Success      	200 {string} string "Shipment deleted successfully"
// @Router 			/ShipmentCRUD/deleteShipment [delete]
func DeleteShipmentHandler(w http.ResponseWriter, r *http.Request) {

	// Extract the id from the URL segment

	id, err := strconv.Atoi(r.URL.Query().Get("shipmentId"))
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get a connection to the database
	db := commons.GetConnection()
	defer db.Close()

	// Start a transaction in the database
	tx, err := db.Begin()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to delete shipment (CONNECTING DB)")
		return
	}
	code, err := deleteShipment(id, tx)
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
	utils.JSONResponse(w, http.StatusOK, "shipment deleted successfully")
}

// private methods

func createShipment(s *models.Shipment, tx *sql.Tx) (int, error) {
	//Check the required fields

	if s.Customer_id == 0 {
		return http.StatusBadRequest, errors.New("customer_id can not be empty")
	}

	if s.Address == "" {
		return http.StatusBadRequest, errors.New("address can not be empty")

	}

	if s.Zip_code == "" { //TODO: use some library for check if we have a valid zip code
		return http.StatusBadRequest, errors.New("zip_code can not be empty")
	}

	// We prepare the query to create a shipment
	stmt, err := tx.Prepare("INSERT INTO Shipment(customer_id, address, city, state, country, zip_code, shipment_date) VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	defer stmt.Close()

	// We execute SQL to insert the shipment
	result, err := stmt.Exec(s.Customer_id, s.Address, s.City, s.State, s.Country, s.Zip_code, s.Shipment_date)
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
		return http.StatusNotModified, errors.New("shipment not created")
	}

	// Get the ID of the newly created user
	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return http.StatusNotModified, err
	}

	// Set the ID of the user struct to the new ID
	s.Shipment_id = uint64(id)

	return http.StatusCreated, nil
}

func updateShipment(s *models.Shipment, tx *sql.Tx) (int, error) {
	//Check the required fields

	if s.Shipment_id == 0 {
		return http.StatusBadRequest, errors.New("shipment_id can not be empty")
	}

	if s.Customer_id == 0 {
		return http.StatusBadRequest, errors.New("customer_id can not be empty")
	}

	if s.Address == "" {
		return http.StatusBadRequest, errors.New("address can not be empty")

	}

	if s.Zip_code == "" { //TODO: use some library for check if we have a valid zip code
		return http.StatusBadRequest, errors.New("zip_code can not be empty")
	}

	if s.Shipment_date == "" { //TODO: use time library for check if we have a valid DATE
		return http.StatusBadRequest, errors.New("shipment_date can not be empty")
	}

	// Se prepara la sentencia SQL para updatear el Shipment
	stmt, err := tx.Prepare("UPDATE shipment SET customer_id = ?, address = ?, city = ?, state = ?, country = ?, zip_code = ?, shipment_date = ? WHERE shipment_id = ?;")
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	defer stmt.Close()

	// Se ejecuta la sentencia SQL para updatear el Shipment
	result, err := stmt.Exec(s.Customer_id, s.Address, s.City, s.State, s.Country, s.Zip_code, s.Shipment_date, s.Shipment_id)
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
		return http.StatusNotModified, errors.New("shipment not updated (NO CHANGES)")
	}

	return http.StatusOK, nil
}

func deleteShipment(id int, tx *sql.Tx) (int, error) {

	// Prepare the SQL statement to delete the shipment
	stmt, err := tx.Prepare("DELETE FROM Shipment WHERE shipment_id = ?")
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to delete shipment (PREPARE QUERY)")
	}
	defer stmt.Close()

	// Execute the SQL statement to delete the shipment
	result, err := stmt.Exec(id)
	if err != nil {
		tx.Rollback()
		return http.StatusNotAcceptable, errors.New("failed to delete shipment (EXECUTE QUERY)")
	}

	// Get the number of rows affected by the SQL statement
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to delete shipment (ROWS AFFECTED)")
	}

	// If no rows were deleted, rollback the transaction
	if rowsAffected == 0 {
		tx.Rollback()
		return http.StatusNotModified, errors.New("failed to delete shipment (NO ROWS AFFECTED)")
	}
	return http.StatusOK, nil
}
