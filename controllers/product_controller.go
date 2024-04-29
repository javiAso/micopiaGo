package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"micopia/commons"
	"micopia/models"
	"micopia/utils"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

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
// @Success      	200 {object} models.Product
// @Router			/ProductCRUD/createProduct [post]
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

func mapCreateProductReqToProduct(c models.CreateProductRequest) models.Product {
	return models.Product{
		Category_id: c.Category_id,
		Description: c.Description,
		Price:       c.Price,
		Stock:       c.Stock,
	}
}
