package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
)

func Register(r *mux.Router) {
	r.Use(responseHeaderMiddleware)
	r.Use(handlers.CompressHandler)

	addV1Routes(r.PathPrefix("/v1").Subrouter())
}

func addV1Routes(r *mux.Router) {
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users", listUsers).Methods("GET")
	r.HandleFunc("/users/{id}", findUser).Methods("GET")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
}

type requester interface {
}

type responser interface {
}

type userCreateRequest struct {
	User *User `json:"user"`
}

type userCreateResponse struct {
	User *User `json:"user"`
}

// update
type userUpdateRequest struct {
	User *User `json:"user"`
}

type userUpdateResponse struct {
	User *User `json:"user"`
}

// find
type userFindResponse struct {
	User *User `json:"user"`
}

// list
type userListResponse struct {
	Users []*User `json:"users"`
}

// error
type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

func encodeRequestBody(payload interface{}) io.Reader {
	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("ERR_ENCODE_REQUEST_BODY	err:%v", err)
		return nil
	}
	return bytes.NewBuffer(body)
}

func decodeResponseBody(bytes []byte, v interface{}) {
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		log.Printf("ERR_DECODE_RESPONSE_BODY	err:%v", err)
		return
	}
}

func decodeRequestBody(body io.ReadCloser, v interface{}) error {
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&v); err != nil {
		log.Printf("Invalid request body	err:%v	body:%v", err, body)
		return fmt.Errorf("Invalid request body")
	}
	defer body.Close()
	return nil
}

func writeErrorResponse(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	res := &errorResponse{
		ErrorMessage: message,
	}
	json.NewEncoder(w).Encode(res)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var p userCreateRequest
	err := decodeRequestBody(r.Body, &p)
	if err != nil {
		writeErrorResponse(w, err.Error())
		return
	}

	if p.User == nil {
		log.Printf("Invalid request payloads	payloads:%v", p)
		writeErrorResponse(w, "Invalid parameter")
		return
	}

	if p.User.Name == "" {
		writeErrorResponse(w, "User name is empty")
		return
	}

	user := &User{
		Id:        uuid.New().String(),
		Name:      p.User.Name,
		CreatedAt: time.Now(),
	}

	repository := UserDatastoreRepository{}
	err = repository.Create(ctx, user)
	if err != nil {
		log.Printf("UserCreateError	err:%v", err)
	}

	res := &userCreateResponse{User: user}
	json.NewEncoder(w).Encode(res)
}

func findUser(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	vars := mux.Vars(r)
	id := vars["id"]

	repository := UserDatastoreRepository{}
	user, err := repository.Find(ctx, id)

	if err != nil {
		log.Printf("FindUser	err:%v", err)
		writeErrorResponse(w, "Can not find user")
		return
	}

	res := userFindResponse{
		User: user,
	}
	json.NewEncoder(w).Encode(res)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	vars := mux.Vars(r)
	id := vars["id"]

	repository := UserDatastoreRepository{}
	err := repository.Delete(ctx, id)
	if err != nil {
		log.Printf("DeleteUser	err:%v", err)
		writeErrorResponse(w, "Can not delete user")
		return
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	vars := mux.Vars(r)
	id := vars["id"]

	var p userUpdateRequest
	err := decodeRequestBody(r.Body, &p)
	if err != nil {
		writeErrorResponse(w, err.Error())
		return
	}

	if p.User == nil {
		log.Printf("Invalid request payloads	payloads:%v", p)
		writeErrorResponse(w, "Invalid parameter")
		return
	}

	repository := UserDatastoreRepository{}
	user, err := repository.Find(ctx, id)

	if err != nil || user == nil {
		log.Printf("FindUser	err:%v", err)
		writeErrorResponse(w, "Can not find user")
		return
	}

	user.Name = p.User.Name

	err = repository.Update(ctx, user)
	if err != nil || user == nil {
		log.Printf("DeleteUser	err:%v", err)
		writeErrorResponse(w, "Can not update user")
		return
	}

	res := userUpdateResponse{
		User: user,
	}
	json.NewEncoder(w).Encode(res)
}

func listUsers(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	repository := UserDatastoreRepository{}
	users, err := repository.List(ctx)
	if err != nil {
		log.Fatalf("ListUser	err:%v", err)
	}

	res := userListResponse{
		Users: users,
	}
	json.NewEncoder(w).Encode(res)
}
