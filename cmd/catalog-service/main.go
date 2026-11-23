package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type product struct {
	ProductID      string `json:"productId"`
	Name           string `json:"name"`
	AvailableStock int    `json:"availableStock"`
}

type api struct {
	products map[string]product
}

func main() {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8081"
	}

	application := api{products: map[string]product{
		"book-123": {ProductID: "book-123", Name: "Go Basics", AvailableStock: 5},
		"pen-456":  {ProductID: "pen-456", Name: "Pen", AvailableStock: 1},
	}}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", application.health)
	mux.HandleFunc("GET /products/{productId}", application.getProduct)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("catalog-service listening on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server: %v", err)
	}
}

func (a api) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a api) getProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("productId")
	item, found := a.products[productID]
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "product not found"})
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("could not encode HTTP response: %v", err)
	}
}
