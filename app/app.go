package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func Start() {

	fmt.Println("Application started...")

	router := mux.NewRouter()

	uh := UploadHandlers{}

	router.HandleFunc("/upload", uh.uploadFiles).Methods(http.MethodPost)

	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":8000", handler))
}
