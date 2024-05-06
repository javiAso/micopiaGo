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
// @Summary: 		Get All Orders
// @Description  	Get All Orders from the database
// @Produce 		application/json
// @Tags			Order
// @Success			200 {object} []models.Order
// @Router			/OrderCRUD/getOrders [get]
func AllOrdersHandler(w http.ResponseWriter, r *http.Request) {

	db := commons.GetConnection()

	var orders []models.Order
	rows, err := db.Query("SELECT o.order_id, o.order_date, o.total_price, o.customer_id, o.payment_id, o.shipment_id FROM `order` o;")
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.Order_id, &order.Order_date, &order.Total_price, &order.Customer_id, &order.Payment_id, &order.Shipment_id); err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// We return Json response:
	utils.JSONResponse(w, http.StatusOK, orders)

}

// CreateTags		godoc
// @Summary: 		Get All Orders
// @Description  	Get All Orders from the database
// @Param			customerId query string true "The customer identifier"
// @Produce 		application/json
// @Tags			Order
// @Success			200 {object} []models.Order
// @Router			/OrderCRUD/getCustomerOrders [get]
func GetCustomerOrdersHandler(w http.ResponseWriter, r *http.Request) {

	// We get "customerId" from URL
	id := r.URL.Query().Get("customerId")
	if id == "" {
		utils.JSONError(w, http.StatusBadRequest, "missing customerId query parameter")
		return
	}

	db := commons.GetConnection()

	var orders []models.Order
	rows, err := db.Query("SELECT o.order_id, o.order_date, o.total_price, o.customer_id, o.payment_id, o.shipment_id FROM `order` o WHERE customer_id=?;", id)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.Order_id, &order.Order_date, &order.Total_price, &order.Customer_id, &order.Payment_id, &order.Shipment_id); err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// We return Json response:
	utils.JSONResponse(w, http.StatusOK, orders)

}

// CreateTags		godoc
// @Summary: 		Get All Orders
// @Description  	Get All Orders from the database
// @Param			paymentId query string true "The payment identifier"
// @Produce 		application/json
// @Tags			Order
// @Success			200 {object} []models.Order
// @Router			/OrderCRUD/getPaymentOrders [get]
func GetPaymentOrdersHandler(w http.ResponseWriter, r *http.Request) {

	// We get "paymentId" from URL
	id := r.URL.Query().Get("paymentId")
	if id == "" {
		utils.JSONError(w, http.StatusBadRequest, "missing paymentId query parameter")
		return
	}

	db := commons.GetConnection()

	var orders []models.Order
	rows, err := db.Query("SELECT o.order_id, o.order_date, o.total_price, o.customer_id, o.payment_id, o.shipment_id FROM `order` o WHERE payment_id=?;", id)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.Order_id, &order.Order_date, &order.Total_price, &order.Customer_id, &order.Payment_id, &order.Shipment_id); err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// We return Json response:
	utils.JSONResponse(w, http.StatusOK, orders)

}

// CreateTags		godoc
// @Summary: 		Get All Orders
// @Description  	Get All Orders from the database
// @Param			shipmentId query string true "The shipment identifier"
// @Produce 		application/json
// @Tags			Order
// @Success			200 {object} []models.Order
// @Router			/OrderCRUD/getShipmentOrders [get]
func GetShipmentOrdersHandler(w http.ResponseWriter, r *http.Request) {

	// We get "shipmentId" from URL
	id := r.URL.Query().Get("shipmentId")
	if id == "" {
		utils.JSONError(w, http.StatusBadRequest, "missing shipmentId query parameter")
		return
	}

	db := commons.GetConnection()

	var orders []models.Order
	rows, err := db.Query("SELECT o.order_id, o.order_date, o.total_price, o.customer_id, o.payment_id, o.shipment_id FROM `order` o WHERE shipment_id=?;", id)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.Order_id, &order.Order_date, &order.Total_price, &order.Customer_id, &order.Payment_id, &order.Shipment_id); err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// We return Json response:
	utils.JSONResponse(w, http.StatusOK, orders)

}

// CreateTags		godoc
// @Summary: 		Get Order by Id
// @Description  	Get Order by Id from the database
// @Param			orderId query string true "The order identifier"
// @Produce 		application/json
// @Tags			Order
// @Success			200 {object} models.Order
// @Router			/OrderCRUD/getOrder [get]
func GetOrderHandler(w http.ResponseWriter, r *http.Request) {

	// We get "orderId" from URL
	id := r.URL.Query().Get("orderId")
	if id == "" {
		utils.JSONError(w, http.StatusBadRequest, "missing orderId query parameter")
		return
	}

	db := commons.GetConnection()

	rows, err := db.Query("SELECT o.order_id, o.order_date, o.total_price, o.customer_id, o.payment_id, o.shipment_id FROM `order` o WHERE order_id=?;", id)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var order models.Order
	for rows.Next() {
		if err := rows.Scan(&order.Order_id, &order.Order_date, &order.Total_price, &order.Customer_id, &order.Payment_id, &order.Shipment_id); err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if err := rows.Err(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// We return Json response:
	utils.JSONResponse(w, http.StatusOK, order)

}

// CreateTags		godoc
// @Summary: 		Create Order
// @Description  	Create Order in the database
// @Param			CreateOrderRequest body models.Order true "The Order to create, order_id is not relevant"
// @Produce 		application/json
// @Tags			Order
// @Success      	201 {object} models.Order
// @Router			/OrderCRUD/createOrder [put]
func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Decodificar el cuerpo de la petición en una estructura CreateOrderRequest
	var createReq models.Order
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed in order Request")
		return
	}

	// Obtener una conexión a la base de datos
	db := commons.GetConnection()
	defer db.Close()

	// Iniciar transacción en la base de datos
	tx, err := db.Begin()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to create order (Internal server error DB.BEGIN)")
		return
	}

	code, err := createOrder(&createReq, tx)
	if err != nil {
		utils.JSONError(w, code, err.Error())
		return
	}

	// Realizar commit de la transacción en la base de datos
	if err := tx.Commit(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to create order (COMMIT FAILED)")
		return
	}

	// Devolver el nuevo Order como JSON
	utils.JSONResponse(w, http.StatusCreated, createReq)
}

// CreateTags		godoc
// @Summary: 		Update Order
// @Description  	Update Order in the database
// @Param			UpdateOrderRequest body models.Order true "The Order to update"
// @Produce 		application/json
// @Tags			Order
// @Success      	200 {object} models.Order
// @Router			/OrderCRUD/updateOrder [post]
func UpdateOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Decodificar el cuerpo de la petición en una estructura Order
	var o models.Order
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
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

	code, err := updateOrder(&o, tx)
	if err != nil {
		utils.JSONError(w, code, err.Error())
		return
	}

	// Realizar commit de la transacción en la base de datos
	if err := tx.Commit(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to update shipment (COMMIT FAILED)")
		return
	}

	// Devolver el nuevo Order como JSON
	utils.JSONResponse(w, http.StatusOK, o)
}

// CreateTags		godoc
// @Summary: 		Delete Order
// @Description  	Delete Order in the database
// @Param			orderId query string true "The Order identifier"
// @Produce 		application/json
// @Tags			Order
// @Success      	200 {string} string "Order deleted successfully"
// @Router 			/OrderCRUD/deleteOrder [delete]
func DeleteOrderHandler(w http.ResponseWriter, r *http.Request) {

	// Extract the id from the URL segment

	id, err := strconv.Atoi(r.URL.Query().Get("orderId"))
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
		utils.JSONError(w, http.StatusInternalServerError, "failed to delete order (CONNECTING DB)")
		return
	}
	code, err := deleteOrder(id, tx)
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
	utils.JSONResponse(w, http.StatusOK, "order deleted successfully")
}

// private methods

func createOrder(o *models.Order, tx *sql.Tx) (int, error) {

	//Check the required fields

	if o.Customer_id == 0 {
		return http.StatusBadRequest, errors.New("customer_id can not be empty")
	}

	if o.Payment_id == 0 {
		return http.StatusBadRequest, errors.New("payment_id can not be empty")
	}

	if o.Shipment_id == 0 {
		return http.StatusBadRequest, errors.New("shipment_id can not be empty")
	}

	if o.Order_date == "" {
		return http.StatusBadRequest, errors.New("order_date can not be empty")
	}

	// We prepare the query to create a order
	stmt, err := tx.Prepare("INSERT INTO `Order`(order_date, total_price, customer_id, payment_id, shipment_id) VALUES(?,?,?,?,?)")
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	defer stmt.Close()

	// We execute SQL to insert the order
	result, err := stmt.Exec(o.Order_date, o.Total_price, o.Customer_id, o.Payment_id, o.Shipment_id)
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
		return http.StatusNotModified, errors.New("order not created")
	}

	// Get the ID of the newly created user
	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return http.StatusNotModified, err
	}

	// Set the ID of the user struct to the new ID
	o.Order_id = uint64(id)

	return http.StatusCreated, nil
}

func updateOrder(o *models.Order, tx *sql.Tx) (int, error) {
	//Check the required fields

	if o.Order_id == 0 {
		return http.StatusBadRequest, errors.New("order_id can not be empty")
	}

	if o.Customer_id == 0 {
		return http.StatusBadRequest, errors.New("customer_id can not be empty")
	}

	if o.Payment_id == 0 {
		return http.StatusBadRequest, errors.New("payment_id can not be empty")
	}

	if o.Shipment_id == 0 {
		return http.StatusBadRequest, errors.New("shipment_id can not be empty")
	}

	if o.Order_date == "" {
		return http.StatusBadRequest, errors.New("order_date can not be empty")
	}

	// Se prepara la sentencia SQL para updatear el Order
	stmt, err := tx.Prepare("UPDATE `order` SET order_date = ?, total_price = ?, customer_id = ?, payment_id = ?, shipment_id = ? WHERE order_id = ?;")
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	defer stmt.Close()

	// Se ejecuta la sentencia SQL para updatear el Order
	result, err := stmt.Exec(o.Order_date, o.Total_price, o.Customer_id, o.Payment_id, o.Shipment_id, o.Order_id)
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
		return http.StatusNotModified, errors.New("order not updated (NO CHANGES)")
	}

	return http.StatusOK, nil
}

func deleteOrder(id int, tx *sql.Tx) (int, error) {

	// Prepare the SQL statement to delete the order
	stmt, err := tx.Prepare("DELETE FROM `Order` WHERE order_id = ?")
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to delete order (PREPARE QUERY)")
	}
	defer stmt.Close()

	// Execute the SQL statement to delete the order
	result, err := stmt.Exec(id)
	if err != nil {
		tx.Rollback()
		return http.StatusNotAcceptable, errors.New("failed to delete order (EXECUTE QUERY)")
	}

	// Get the number of rows affected by the SQL statement
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to delete order (ROWS AFFECTED)")
	}

	// If no rows were deleted, rollback the transaction
	if rowsAffected == 0 {
		tx.Rollback()
		return http.StatusNotModified, errors.New("failed to delete order (NO ROWS AFFECTED)")
	}
	return http.StatusOK, nil
}
