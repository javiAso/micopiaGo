package commons

import "net/http"

func SendResponse(writer http.ResponseWriter, status int, data []byte) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Credentials", "true")
	writer.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	writer.Write(data)

}

func SendError(writer http.ResponseWriter, status int) {
	data := []byte(`{}`)
	writer.Header().Set("Content-Type", "application/json")

	writer.WriteHeader(status)
	writer.Write(data)
}

func SendResponseInsert(writer http.ResponseWriter, status int) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Credentials", "true")
	writer.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)

}
