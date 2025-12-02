package main

import (
	"log"
	"net/http"
)

func main() {

	router := http.NewServeMux()
	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Apı get request received")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello  ı am api test service "))
	})
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Apı get request received ")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello  ı am api test service "))
	})
	log.Println("Starting server on :8082")
	if err := http.ListenAndServe(":8082", router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
