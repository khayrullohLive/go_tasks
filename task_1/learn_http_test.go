package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHome(t *testing.T) {
	w := httptest.NewRecorder()
	Home(w, nil)
	districtCode := http.StatusOK
	if w.Code != districtCode {
		t.Errorf("Home() returned wrong status code: got %v want %v", w.Code, districtCode)
	}
	expected := []byte("My Home")
	if !bytes.Equal(w.Body.Bytes(), expected) {
		t.Errorf("Home() returned wrong body: got %v want %v", w.Body.Bytes(), expected)
	}
}
func TestTodos(t *testing.T) {
	w := httptest.NewRecorder()
	Todos(w, nil)
	districtCode := http.StatusOK
	if w.Code != districtCode {
		t.Errorf("Todos() returned wrong status code: got %v want %v", w.Code, districtCode)
	}
	expected := []byte("[{\"title\":\"bay book\"}]")
	if !bytes.Equal(w.Body.Bytes(), expected) {
		t.Errorf("Todos() returned wrong body: got %v want %v", w.Body.Bytes(), expected)
	}

}
