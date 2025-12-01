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
		w.Write([]byte("Hello  ı am api gateway service "))
	})

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
