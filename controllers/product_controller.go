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
) //TODO: implement rollbacks

// CreateTags		godoc
// @Summary: 		Get All Products
// @Description  	Get All Products from the database
// @Produce 		application/json
// @Tags			Product
// @Success			200 {object} models.Products
// @Router			/ProductCRUD/getProducts [get]
func AllProductsHandler(w http.ResponseWriter, r *http.Request) {

	db := commons.GetConnection()

	var Products []models.Product
	rows, err := db.Query("SELECT p.product_id, p.description, p.price, p.stock, p.category_id FROM product p;")
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.Product_id, &product.Description, &product.Price, &product.Stock, &product.Category_id); err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		Products = append(Products, product)
	}

	if err := rows.Err(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var productList models.Products
	productList.ProductList = Products

	// Convertir el resultado en formato JSON y enviar la respuesta
	utils.JSONResponse(w, http.StatusOK, productList)

}

// CreateTags		godoc
// @Summary: 		Create Product
// @Description  	Create Product in the database
// @Param			CreateProductRequest body models.CreateProductRequest true "The Product to create"
// @Produce 		application/json
// @Tags			Product
// @Success      	201 {object} models.Product
// @Router			/ProductCRUD/createProduct [put]
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Decodificar el cuerpo de la petición en una estructura CreateProductRequest
	var createReq models.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Failed in Product Request")
		return
	}

	// Obtener una conexión a la base de datos
	db := commons.GetConnection()
	defer db.Close()

	// Iniciar transacción en la base de datos
	tx, err := db.Begin()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Failed to create Product (Internal server error)")
		return
	}

	product := mapCreateProductReqToProduct(createReq)

	code, err := createProduct(&product, tx)
	if err != nil {
		utils.JSONError(w, code, err.Error())
		return
	}

	// Realizar commit de la transacción en la base de datos
	if err := tx.Commit(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Failed to create Products (COMMIT FAILED)")
		return
	}

	// Devolver el nuevo producto como JSON
	utils.JSONResponse(w, http.StatusCreated, product)
}

// CreateTags		godoc
// @Summary: 		Update Product
// @Description  	Update Product in the database
// @Param			UpdateProductRequest body models.Product true "The Product to update"
// @Produce 		application/json
// @Tags			Product
// @Success      	200 {object} models.Product
// @Router			/ProductCRUD/updateProduct [post]
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Decodificar el cuerpo de la petición en una estructura Product
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Failed in Product Request")
		return
	}

	// Obtener una conexión a la base de datos
	db := commons.GetConnection()
	defer db.Close()

	// Iniciar transacción en la base de datos
	tx, err := db.Begin()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Failed to create Product (Internal server error)")
		return
	}

	code, err := updateProduct(&product, tx)
	if err != nil {
		utils.JSONError(w, code, err.Error())
		return
	}

	// Realizar commit de la transacción en la base de datos
	if err := tx.Commit(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Failed to update Products (COMMIT FAILED)")
		return
	}

	// Devolver el nuevo producto como JSON
	utils.JSONResponse(w, http.StatusOK, product)
}

// CreateTags		godoc
// @Summary: 		Delete Product
// @Description  	Delete Product in the database
// @Param			productId query string true "The product identifier"
// @Produce 		application/json
// @Tags			Product
// @Success      	200 {string} string "deleted"
// @Router 			/ProductCRUD/deleteProduct [delete]
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Get a connection to the database
	db := commons.GetConnection()
	defer db.Close()

	// Extract the id from the URL segment

	id, err := strconv.Atoi(r.URL.Query().Get("productId"))
	if err != nil {
		println(err)
		utils.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	// Start a transaction in the database
	tx, err := db.Begin()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Failed to delete Product (CONNECTING DB)")
		return
	}
	code, err := deleteProduct(id, tx)
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
	utils.JSONResponse(w, http.StatusOK, "Product deleted successfully")
}

//Private methods:

func createProduct(p *models.Product, tx *sql.Tx) (int, error) {
	//Check the required fields
	if p.Description == "" {
		return http.StatusBadRequest, errors.New("description can not be empty")
	}
	if p.Category_id == 0 {
		return http.StatusBadRequest, errors.New("category can not be empty")
	}

	// Se prepara la sentencia SQL para insertar la hora
	stmt, err := tx.Prepare("INSERT INTO Product(description, price, stock, category_id) VALUES(?, ?, ?, ?)")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer stmt.Close()

	// Se ejecuta la sentencia SQL para insertar el user
	result, err := stmt.Exec(p.Description, p.Price, p.Stock, p.Category_id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Se obtiene la cantidad de filas afectadas por la sentencia SQL
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// Si no se ha insertado ninguna fila, se hace rollback de la transacción
	if rowsAffected == 0 {
		return http.StatusInternalServerError, errors.New("product not created")
	}

	// Get the ID of the newly created user
	id, err := result.LastInsertId()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Set the ID of the user struct to the new ID
	p.Product_id = uint64(id)

	return http.StatusCreated, nil
}

func updateProduct(p *models.Product, tx *sql.Tx) (int, error) {
	//Check the required fields

	if p.Product_id == 0 {
		return http.StatusBadRequest, errors.New("product_id can not be empty")
	}

	if p.Description == "" {
		return http.StatusBadRequest, errors.New("description can not be empty")
	}
	if p.Category_id == 0 {
		return http.StatusBadRequest, errors.New("category can not be empty")
	}

	// Se prepara la sentencia SQL para insertar el producto
	stmt, err := tx.Prepare("UPDATE Product SET description = ?, price = ?, stock = ?, category_id = ? WHERE product_id = ?")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer stmt.Close()

	// Se ejecuta la sentencia SQL para updatear el producto
	result, err := stmt.Exec(p.Description, p.Price, p.Stock, p.Category_id, p.Product_id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Se obtiene la cantidad de filas afectadas por la sentencia SQL
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// Si no se ha updateado ninguna fila, se hace rollback de la transacción
	if rowsAffected == 0 {
		return http.StatusNotModified, errors.New("product not updated (NO CHANGES)")
	}

	return http.StatusOK, nil
}

func deleteProduct(id int, tx *sql.Tx) (int, error) {

	// Prepare the SQL statement to delete the product
	stmt, err := tx.Prepare("DELETE FROM Product WHERE product_id = ?")
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to delete product (PREPARE QUERY)")
	}
	defer stmt.Close()

	// Execute the SQL statement to delete the product
	result, err := stmt.Exec(id)
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to delete product (EXECUTE QUERY)")
	}

	// Get the number of rows affected by the SQL statement
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to delete product (ROWS AFFECTED)")
	}

	// If no rows were deleted, rollback the transaction
	if rowsAffected == 0 {
		tx.Rollback()
		return http.StatusNotModified, errors.New("failed to delete product (NO ROWS AFFECTED)")
	}
	return http.StatusOK, nil
}

func mapCreateProductReqToProduct(c models.CreateProductRequest) models.Product {
	return models.Product{
		Category_id: c.Category_id,
		Description: c.Description,
		Price:       c.Price,
		Stock:       c.Stock,
	}
}
