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
		w.Write([]byte("Hello  ı am api user service "))
	})
	router.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Apı get request received")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello  ı am api user service -profile"))
	})
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Apı get request received ")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello  ı am api user service "))
	})
	log.Println("Starting server on :8083")
	if err := http.ListenAndServe(":8083", router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
