package main

import (
  "fmt"
  "net/http"
)

func main() {

  // create mux (an http request multiplexer) to handle incoming requests
  // and redirect the requests to the correct handler function by 
  // pattern matching
  mux := http.NewServeMux()

  // create handler function for requests to "/"
  mux.HandleFunc("/", handleRoot)

  // start server and print start message
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
