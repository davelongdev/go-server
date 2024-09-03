package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

// create struct to  store user data
type User struct {
  Name string `json:"name"`
}

// create simple key value store database using a map
var userKV = make(map[int]User)

// create read and write mutex to enable threadsafe alteration of 
// the kv store by locking and unlocking it
var kvMutex sync.RWMutex

func main() {

  // create mux (an http request multiplexer) to handle incoming requests
  // and redirect the requests to the correct handler function by 
  // pattern matching
  mux := http.NewServeMux()

  // create handler function for requests to "/"
  mux.HandleFunc("/", handleRoot)

  // create handler function for POST requests to "/users"
  mux.HandleFunc("POST /users", createUser)

  // create handler function for GET requests to "/users"
  mux.HandleFunc("GET /users/{id}", getUser)

  // create handler function for DELETE requests to "/users"
  mux.HandleFunc("DELETE /users/{id}", deleteUser)

  // start server and print start message
  fmt.Println("Server listening to :8080")
  http.ListenAndServe(":8080", mux)

}

func handleRoot(

  // w is responsible for constructing an http response
  w http.ResponseWriter,

  // r contains an http request
  r *http.Request,

) {

  // print success message
  fmt.Fprintf(w, "Root directory successfully accessed")
}

func getUser(
  w http.ResponseWriter,
  r *http.Request,
) {

  // get user id
  id, err := strconv.Atoi(r.PathValue("id"))

  // handle error 
  if err != nil {
    http.Error(
      w,
      err.Error(),
      http.StatusBadGateway,
    )
    return
  }

  // get user from userKV database safely using mutex locking and unlocking
  kvMutex.RLock()
  user, ok := userKV[id]
  kvMutex.RUnlock()

  // handle error
  if !ok {
    http.Error(
      w,
      "user not found",
      http.StatusNotFound,
    )
    return
  }
  
  // set content type header for http response
  w.Header().Set("Content-Type","application/json")

  // get a json representation of user data
  j, err := json.Marshal(user)

  // handle error
  if err != nil {
    http.Error(
      w,
      err.Error(),
      http.StatusInternalServerError,
    )
    return
  }

  // write status header
  w.WriteHeader(http.StatusOK)

  // write body containing json
  w.Write(j)
}

func createUser(
  w http.ResponseWriter,
  r *http.Request,
) {

  // create user variable
  var user User

  // decode json and pass data to user variable
  err := json.NewDecoder(r.Body).Decode(&user)

  // handle error
  if err != nil {
    http.Error(
      w, 
      err.Error(), 
      http.StatusBadRequest,
    )
    return 
  }

  // handle error (if request doesn't send a  namee)
  if user.Name == "" {
    http.Error(
      w, 
      "name is required", 
      http.StatusBadRequest,
    )
    return
  }

  // add user to database
  kvMutex.Lock()
  userKV[len(userKV) + 1] = user
  kvMutex.Unlock()

  // write http status header 
  w.WriteHeader(http.StatusNoContent)
}

func deleteUser(
  w http.ResponseWriter,
  r *http.Request,
) {

  // get user id
  id, err := strconv.Atoi(r.PathValue("id"))

  // handle error
  if err != nil {
    http.Error(
      w,
      err.Error(),
      http.StatusBadRequest,
    )
    return
  }

  // handle error
  if _, ok := userKV[id]; !ok {
    http.Error(
     w,
      "user not found",
      http.StatusBadRequest,
    )
    return
  }

  // delete user from database
  kvMutex.Lock()
  delete(userKV, id)
  kvMutex.Unlock()

  // write status header for http response
  w.WriteHeader(http.StatusNoContent)
}
