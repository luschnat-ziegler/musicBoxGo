package app

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type UploadHandlers struct {}

func (uh *UploadHandlers) uploadFiles (w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(200000) // grab the multipart form
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	formData := r.MultipartForm
	files := formData.File["multiplefiles"]

	for _, file := range files {
		fmt.Println(file.Filename)
	}

	response := struct {
		Success bool
	}{true}

	writeResponse(w, http.StatusOK, response)

}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}