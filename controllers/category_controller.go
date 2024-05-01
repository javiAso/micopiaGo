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
// @Summary: 		Get All Categories
// @Description  	Get All Categories from the database
// @Produce 		application/json
// @Tags			Category
// @Success			200 {object} models.Categories
// @Router			/CategoryCRUD/getCategories [get]
func AllCategorysHandler(w http.ResponseWriter, r *http.Request) {

	db := commons.GetConnection()

	rows, err := db.Query("SELECT c.category_id, c.name FROM Category c;")
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.Category_id, &c.Name); err != nil {
			utils.JSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var CategoryList models.Categories
	CategoryList.CategoryList = categories

	// Convertir el resultado en formato JSON y enviar la respuesta
	utils.JSONResponse(w, http.StatusOK, CategoryList)

}

// CreateTags		godoc
// @Summary: 		Get Category
// @Description  	Get Category from the database by id
// @Param			categoryId query string true "The Category identifier"
// @Produce 		application/json
// @Tags			Category
// @Success			200 {object} models.Category
// @Router			/CategoryCRUD/getCategory [get]
func GetCategoryHandler(w http.ResponseWriter, r *http.Request) {

	// Obtener el ID de la categoría del parámetro de consulta "categoryId"
	id := r.URL.Query().Get("categoryId")
	if id == "" {
		utils.JSONError(w, http.StatusBadRequest, "missing categoryId query parameter")
		return
	}

	// Obtener una conexión a la base de datos
	db := commons.GetConnection()
	defer db.Close()

	// Realizar la consulta SQL para obtener las horas de la categoría
	rows, err := db.Query("SELECT c.category_id, c.name FROM Category c WHERE c.category_id = ?;", id)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to get category (DB QUERY ERROR)")
		return
	}
	defer rows.Close()

	//We scan the Data
	var c models.Category
	for rows.Next() {
		if err := rows.Scan(&c.Category_id, &c.Name); err != nil {
			utils.JSONError(w, http.StatusInternalServerError, "failed to get category ( ROWS SCAN ERROR)")
			return
		}
	}

	if err := rows.Err(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to get category (ROWS ERROR)")
		return
	}

	// Devolver la category como JSON
	utils.JSONResponse(w, http.StatusOK, c)

}

// CreateTags		godoc
// @Summary: 		Create Category
// @Description  	Create Category in the database
// @Param			CreateCategoryRequest body models.CreateCategoryRequest true "The Category to create"
// @Produce 		application/json
// @Tags			Category
// @Success      	201 {object} models.Category
// @Router			/CategoryCRUD/createCategory [put]
func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Decodificar el cuerpo de la petición en una estructura CreateCategoryRequest
	var createReq models.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed in category Request")
		return
	}

	// Obtener una conexión a la base de datos
	db := commons.GetConnection()
	defer db.Close()

	// Iniciar transacción en la base de datos
	tx, err := db.Begin()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to create category (Internal server error DB.BEGIN)")
		return
	}

	category := new(models.Category)
	category.Name = createReq.Name

	code, err := createCategory(category, tx)
	if err != nil {
		utils.JSONError(w, code, err.Error())
		return
	}

	// Realizar commit de la transacción en la base de datos
	if err := tx.Commit(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to create category (COMMIT FAILED)")
		return
	}

	// Devolver el nuevo Category como JSON
	utils.JSONResponse(w, http.StatusCreated, category)
}

// CreateTags		godoc
// @Summary: 		Update Category
// @Description  	Update Category in the database
// @Param			UpdateCategoryRequest body models.Category true "The Category to update"
// @Produce 		application/json
// @Tags			Category
// @Success      	200 {object} models.Category
// @Router			/CategoryCRUD/updateCategory [post]
func UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Decodificar el cuerpo de la petición en una estructura Category
	var c models.Category
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed in category request")
		return
	}

	// Obtener una conexión a la base de datos
	db := commons.GetConnection()
	defer db.Close()

	// Iniciar transacción en la base de datos
	tx, err := db.Begin()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to update category (Internal server error DD.BEGIN)")
		return
	}

	code, err := updateCategory(&c, tx)
	if err != nil {
		utils.JSONError(w, code, err.Error())
		return
	}

	// Realizar commit de la transacción en la base de datos
	if err := tx.Commit(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to update category (COMMIT FAILED)")
		return
	}

	// Devolver el nuevo Category como JSON
	utils.JSONResponse(w, http.StatusOK, c)
}

// CreateTags		godoc
// @Summary: 		Delete Category
// @Description  	Delete Category in the database
// @Param			categoryId query string true "The Category identifier"
// @Produce 		application/json
// @Tags			Category
// @Success      	200 {string} string "Category deleted successfully"
// @Router 			/CategoryCRUD/deleteCategory [delete]
func DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// Get a connection to the database
	db := commons.GetConnection()
	defer db.Close()

	// Extract the id from the URL segment

	id, err := strconv.Atoi(r.URL.Query().Get("categoryId"))
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	// Start a transaction in the database
	tx, err := db.Begin()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to delete category (CONNECTING DB)")
		return
	}
	code, err := deleteCategory(id, tx)
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
	utils.JSONResponse(w, http.StatusOK, "category deleted successfully")
}

// private methods

func createCategory(c *models.Category, tx *sql.Tx) (int, error) {
	//Check the required fields

	if c.Name == "" {
		return http.StatusBadRequest, errors.New("name can not be empty")
	}

	// We prepare the query to create a user
	stmt, err := tx.Prepare("INSERT INTO Category(name) VALUES(?)")
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	defer stmt.Close()

	// Se ejecuta la sentencia SQL para insertar el user
	result, err := stmt.Exec(c.Name)
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
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
		return http.StatusNotModified, errors.New("category not created")
	}

	// Get the ID of the newly created user
	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return http.StatusNotModified, err
	}

	// Set the ID of the user struct to the new ID
	c.Category_id = uint64(id)

	return http.StatusCreated, nil
}

func updateCategory(c *models.Category, tx *sql.Tx) (int, error) {
	//Check the required fields

	if c.Category_id == 0 {
		return http.StatusBadRequest, errors.New("category_id can not be empty")
	}

	if c.Name == "" {
		return http.StatusBadRequest, errors.New("name can not be empty")
	}

	// Se prepara la sentencia SQL para updatear el Category
	stmt, err := tx.Prepare("UPDATE Category SET name = ? WHERE category_id = ?")
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	defer stmt.Close()

	// Se ejecuta la sentencia SQL para updatear el Category
	result, err := stmt.Exec(c.Name, c.Category_id)
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
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
		return http.StatusNotModified, errors.New("category not updated (NO CHANGES)")
	}

	return http.StatusOK, nil
}

func deleteCategory(id int, tx *sql.Tx) (int, error) {

	// Prepare the SQL statement to delete the category
	stmt, err := tx.Prepare("DELETE FROM Category WHERE category_id = ?")
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to delete category (PREPARE QUERY)")
	}
	defer stmt.Close()

	// Execute the SQL statement to delete the category
	result, err := stmt.Exec(id)
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to delete category (EXECUTE QUERY)")
	}

	// Get the number of rows affected by the SQL statement
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, errors.New("failed to delete category (ROWS AFFECTED)")
	}

	// If no rows were deleted, rollback the transaction
	if rowsAffected == 0 {
		tx.Rollback()
		return http.StatusNotModified, errors.New("failed to delete category (NO ROWS AFFECTED)")
	}
	return http.StatusOK, nil
}
