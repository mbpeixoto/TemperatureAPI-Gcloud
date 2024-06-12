package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

type MockTransport struct {
	RoundTripFunc func(req *http.Request) *http.Response
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req), nil
}

func TestHandleTemperatura_Success(t *testing.T) {

	req, err := http.NewRequest("GET", "/temperatura/71503503", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/temperatura/{cep}", HandleTemperatura)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code incorreto: recebeu %v esperava %v", status, http.StatusOK)
	}

}

func TestHandleTemperatura_InvalidCEP(t *testing.T) {
	req, err := http.NewRequest("GET", "/temperatura/123", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/temperatura/{cep}", HandleTemperatura)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("Status code incorreto: recebeu %v esperava %v", status, http.StatusUnprocessableEntity)
	}

	expected := "invalid zipcode\n"
	if rr.Body.String() != expected {
		t.Errorf("Resposta body incorreta: recebeu %v esperava %v", rr.Body.String(), expected)
	}
}

func TestHandleTemperatura_CEPNotFound(t *testing.T) {

	req, err := http.NewRequest("GET", "/temperatura/87654321", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/temperatura/{cep}", HandleTemperatura)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Status code incorreto: recebeu %v esperava %v", status, http.StatusNotFound)
	}

	expected := "can not find zipcode\n"
	if rr.Body.String() != expected {
		t.Errorf("Resposta body incorreta: recebeu %v esperava %v", rr.Body.String(), expected)
	}
}
