package controllers

import (
	"micopia/commons"
	"micopia/models"
	"micopia/utils"
	"net/http"
)

// CreateTags		godoc
// @Summary: 		Get All Categories
// @Description  	Get All Categories from the database
// @Produce 		application/json
// @Tags			Category
// @Success			200 {object} models.Categories
// @Router			/categoryCRUD/getCategories [get]
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
// @Router			/categoryCRUD/getCategory [get]
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
