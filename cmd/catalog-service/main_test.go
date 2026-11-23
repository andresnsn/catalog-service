package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProductReturnsCatalogEntry(t *testing.T) {
	application := api{products: map[string]product{
		"book-123": {ProductID: "book-123", Name: "Go Basics", AvailableStock: 5},
	}}
	request := httptest.NewRequest(http.MethodGet, "/products/book-123", nil)
	request.SetPathValue("productId", "book-123")
	response := httptest.NewRecorder()

	application.getProduct(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if got := response.Body.String(); got != "{\"productId\":\"book-123\",\"name\":\"Go Basics\",\"availableStock\":5}\n" {
		t.Fatalf("body = %q", got)
	}
}

func TestGetProductReturnsNotFoundForUnknownID(t *testing.T) {
	application := api{products: map[string]product{}}
	request := httptest.NewRequest(http.MethodGet, "/products/missing", nil)
	request.SetPathValue("productId", "missing")
	response := httptest.NewRecorder()

	application.getProduct(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNotFound)
	}
}
