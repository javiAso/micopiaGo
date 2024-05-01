package controllers

import (
	"micopia/commons"
	"micopia/models"
	"micopia/utils"
	"net/http"
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
