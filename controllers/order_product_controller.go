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
// @Summary: 		Get All Orders Products
// @Description  	Get All Orders Products from the database
// @Produce 		application/json
// @Tags			OrderProduct
// @Success			200 {object} []models.OrderProduct
// @Router			/OrderProductCRUD/getOrdersProducts [get]
func AllOrdersProductsHandler(w http.ResponseWriter, r *http.Request) {

	db := commons.GetConnection()

	var ordersProducts []models.OrderProduct
	rows, err := db.Query("SELECT o.order_id, o.product_id, o.quantity, o.price FROM order_product o;")
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var oP models.OrderProduct
		if err := rows.Scan(&oP.Order_id, &oP.Product_id, &oP.Quantity, &oP.Price); err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		ordersProducts = append(ordersProducts, oP)
	}

	if err := rows.Err(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// We transform the data into a Json response
	utils.JSONResponse(w, http.StatusOK, ordersProducts)

}

// CreateTags		godoc
// @Summary: 		Get All Order Products
// @Description  	Get All Order Products from the database
// @Param			orderId query string true "The order identifier"
// @Produce 		application/json
// @Tags			OrderProduct
// @Success			200 {object} []models.OrderProduct
// @Router			/OrderProductCRUD/getOrderProducts [get]
func GetOrderProductsHandler(w http.ResponseWriter, r *http.Request) {
	// We get "orderId" from URL
	id := r.URL.Query().Get("orderId")
	if id == "" {
		utils.JSONError(w, http.StatusBadRequest, "missing orderId query parameter")
		return
	}

	db := commons.GetConnection()

	var ordersProducts []models.OrderProduct
	rows, err := db.Query("SELECT o.order_id, o.product_id, o.quantity, o.price FROM order_product o WHERE order_id = ?;", id)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var oP models.OrderProduct
		if err := rows.Scan(&oP.Order_id, &oP.Product_id, &oP.Quantity, &oP.Price); err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		ordersProducts = append(ordersProducts, oP)
	}

	if err := rows.Err(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// We transform the data into a Json response
	utils.JSONResponse(w, http.StatusOK, ordersProducts)

}

// CreateTags		godoc
// @Summary: 		Get All Order Products
// @Description  	Get All Order Products from the database
// @Param			productId query string true "The product identifier"
// @Produce 		application/json
// @Tags			OrderProduct
// @Success			200 {object} []models.OrderProduct
// @Router			/OrderProductCRUD/getProductOrders [get]
func GetProductOrdersHandler(w http.ResponseWriter, r *http.Request) {
	// We get "productId" from URL
	id := r.URL.Query().Get("productId")
	if id == "" {
		utils.JSONError(w, http.StatusBadRequest, "missing productId query parameter")
		return
	}

	db := commons.GetConnection()

	var ordersProducts []models.OrderProduct
	rows, err := db.Query("SELECT o.order_id, o.product_id, o.quantity, o.price FROM order_product o WHERE product_id = ?;", id)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var oP models.OrderProduct
		if err := rows.Scan(&oP.Order_id, &oP.Product_id, &oP.Quantity, &oP.Price); err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		ordersProducts = append(ordersProducts, oP)
	}

	if err := rows.Err(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// We transform the data into a Json response
	utils.JSONResponse(w, http.StatusOK, ordersProducts)

}

// CreateTags		godoc
// @Summary: 		Create OrderProduct
// @Description  	Create OrderProduct in the database
// @Param			CreateOrderProductRequest body models.OrderProduct true "The Order Product to create"
// @Produce 		application/json
// @Tags			OrderProduct
// @Success      	201 {object} models.OrderProduct
// @Router			/OrderProductCRUD/createOrderProduct [put]
func CreateOrderProductHandler(w http.ResponseWriter, r *http.Request) {
	// Decodificar el cuerpo de la petición en una estructura CreateOrderRequest
	var createReq models.OrderProduct
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed in order product Request")
		return
	}

	// Obtener una conexión a la base de datos
	db := commons.GetConnection()
	defer db.Close()

	// Iniciar transacción en la base de datos
	tx, err := db.Begin()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to create order product (Internal server error DB.BEGIN)")
		return
	}

	code, err := createOrderProduct(&createReq, tx)
	if err != nil {
		utils.JSONError(w, code, err.Error())
		return
	}

	// Realizar commit de la transacción en la base de datos
	if err := tx.Commit(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to create order product (COMMIT FAILED)")
		return
	}

	// Devolver el nuevo Order como JSON
	utils.JSONResponse(w, http.StatusCreated, createReq)
}

// CreateTags		godoc
// @Summary: 		Update OrderProduct
// @Description  	Update OrderProduct in the database
// @Param			UpdateOrderProductRequest body models.OrderProduct true "The Order Product to update"
// @Produce 		application/json
// @Tags			OrderProduct
// @Success      	200 {object} models.OrderProduct
// @Router			/OrderProductCRUD/updateOrderProduct [post]
func UpdateOrderProductHandler(w http.ResponseWriter, r *http.Request) {
	// Decodificar el cuerpo de la petición en una estructura Order
	var o models.OrderProduct
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

	code, err := updateOrderProduct(&o, tx)
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
// @Summary: 		Delete OrderProduct
// @Description  	Delete OrderProduct in the database
// @Param			orderId query string true "The OrderProduct identifier"
// @Param			productId query string true "The OrderProduct identifier"
// @Produce 		application/json
// @Tags			OrderProduct
// @Success      	200 {string} string "OrderProduct deleted successfully"
// @Router 			/OrderProductCRUD/deleteOrderProduct [delete]
func DeleteOrderProductHandler(w http.ResponseWriter, r *http.Request) {
	// Get a connection to the database
	db := commons.GetConnection()
	defer db.Close()

	// Extract the id from the URL segment

	oId, err := strconv.Atoi(r.URL.Query().Get("orderId"))
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Extract the id from the URL segment

	pId, err := strconv.Atoi(r.URL.Query().Get("productId"))
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	// Start a transaction in the database
	tx, err := db.Begin()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to delete order product (CONNECTING DB)")
		return
	}
	code, err := deleteOrderProduct(oId, pId, tx)
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
	utils.JSONResponse(w, http.StatusOK, "order product deleted successfully")
}

// private methods

func createOrderProduct(o *models.OrderProduct, tx *sql.Tx) (int, error) {

	//Check the required fields

	if o.Order_id == 0 {
		return http.StatusBadRequest, errors.New("order_id can not be empty")
	}

	if o.Product_id == 0 {
		return http.StatusBadRequest, errors.New("product_id can not be empty")
	}

	// We prepare the query to create a order
	stmt, err := tx.Prepare("INSERT INTO `order_product`(order_id, product_id, price, quantity) VALUES(?,?,?,?)")
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	defer stmt.Close()

	// We execute SQL to insert the order
	result, err := stmt.Exec(o.Order_id, o.Product_id, o.Price, o.Quantity)
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
		return http.StatusNotModified, errors.New("order product not created")
	}

	return http.StatusCreated, nil
}

func updateOrderProduct(o *models.OrderProduct, tx *sql.Tx) (int, error) {
	//Check the required fields

	if o.Order_id == 0 {
		return http.StatusBadRequest, errors.New("order_id can not be empty")
	}

	if o.Product_id == 0 {
		return http.StatusBadRequest, errors.New("product_id can not be empty")
	}

	// We prepare SQL to update order product
	stmt, err := tx.Prepare("UPDATE `order_product` SET price = ?, quantity = ? WHERE order_id = ? AND product_id = ?;")
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	defer stmt.Close()

	// Se ejecuta la sentencia SQL para updatear el Order
	result, err := stmt.Exec(o.Price, o.Quantity, o.Order_id, o.Product_id)
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
		return http.StatusNotModified, errors.New("order product not updated (NO CHANGES)")
	}

	return http.StatusOK, nil
}

func deleteOrderProduct(oId int, pId int, tx *sql.Tx) (int, error) {

	// Prepare the SQL statement to delete the order
	stmt, err := tx.Prepare("DELETE FROM `order_product` WHERE order_id = ? AND product_id = ?")
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to delete order product (PREPARE QUERY)")
	}
	defer stmt.Close()

	// Execute the SQL statement to delete the order
	result, err := stmt.Exec(oId, pId)
	if err != nil {
		tx.Rollback()
		return http.StatusNotAcceptable, errors.New("failed to delete order product (EXECUTE QUERY)")
	}

	// Get the number of rows affected by the SQL statement
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to delete order product (ROWS AFFECTED)")
	}

	// If no rows were deleted, rollback the transaction
	if rowsAffected == 0 {
		tx.Rollback()
		return http.StatusNotModified, errors.New("failed to delete order product (NO ROWS AFFECTED)")
	}
	return http.StatusOK, nil
}
