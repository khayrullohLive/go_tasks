package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
)

func main() {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", Home)
	serveMux.HandleFunc("/todos/", Todos)
	serveMux.HandleFunc("/books/", GetBooks)
	err := http.ListenAndServe(":8080", serveMux)
	if err != nil {
		return
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	response, err := w.Write([]byte("My Home"))
	if err != nil {
		return
	}
	switch r.Method {
	case http.MethodGet:
		fmt.Println("MethodGet")
		return
	case http.MethodPost:
		fmt.Println("MethodPost")
		return
	}
	fmt.Println(response)
}
func Todos(w http.ResponseWriter, r *http.Request) {
	response, err := w.Write([]byte("[{\"title\":\"bay book\"}]"))
	if err != nil {
		return
	}
	fmt.Println(response)
}
func GetBooks(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	names := params["name"]
	var output bytes.Buffer
	output.WriteString("Hello ")
	output.WriteString(names[0])
	_, err := w.Write(output.Bytes())
	if err != nil {
		slog.Error("Error writing to response writer")
		return
	}

}
