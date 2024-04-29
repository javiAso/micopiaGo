package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func JSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errResp := ErrorResponse{statusCode, message}
	json.NewEncoder(w).Encode(errResp)
}

// JSONResponse envía una respuesta JSON con el código de estado HTTP y los datos proporcionados.
func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsonData, err := json.Marshal(data)
	if err != nil {
		// Si hubo un error al codificar el JSON, envía una respuesta de error
		JSONError(w, http.StatusInternalServerError, "Error al codificar la respuesta JSON")
		return
	}
	w.Write(jsonData)
}
