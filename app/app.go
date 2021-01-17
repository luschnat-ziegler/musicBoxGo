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
	mh := MpdHandlers{}

	router.HandleFunc("/upload", uh.uploadFiles).Methods(http.MethodPost)
	router.HandleFunc("/mpdupdate" , mh.updateDB).Methods(http.MethodGet)

	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":8000", handler))
}
