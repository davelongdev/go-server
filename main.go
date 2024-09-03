package main

import (
  "fmt"
  "net/http"
)

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", handleRoot)

  fmt.Println("Server listening to :8080")
  http.ListenAndServe(":8080", mux)
}

func handleRoot(

  // ResponseWriter is responsible for constructing an http response
  w http.ResponseWriter,

  // contains http request
  r *http.Request,
) {
  fmt.Fprintf(w, "Hello World")
}
