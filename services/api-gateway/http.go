package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/shared/contracts"
)

type APIError struct {
	Error string
}

type APIResponse struct {
	Message string `json:"message"`
}

type apiFunc func (w http.ResponseWriter, r *http.Request) error




// Routes
func home(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusAccepted, APIResponse{Message: "Hello World from the API GateWay!"})
}

func handleTripPreview(w http.ResponseWriter, r *http.Request) error{
	var reqBody previewTripRequest
	
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil{
		return WriteJSON(w, http.StatusBadRequest,"failed to parse JSON data")
	}
	
	defer r.Body.Close()
	
	// validation
	if reqBody.UserID == ""{
		return WriteJSON(w,http.StatusBadRequest,"user ID is required")
	}
	
	requestMarshalBody,err := json.Marshal(reqBody)
	if err != nil{
		log.Println("unable to parse the data")
	}
	
	resp, err := http.Post("http://trip-service:8083/preview","application/json",bytes.NewBuffer(requestMarshalBody))
	
	if err != nil{
		log.Println("unable to send a request", err)
	}
	
	defer resp.Body.Close()
	
	var respBody any
	
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil{
		return WriteJSON(w, http.StatusBadRequest,"failed to parse JSON data from trip service")
	}
	// log.Println("Response from the trip :", resp)
	
	
	response := contracts.APIResponse{Data: respBody}

	
	return WriteJSON(w, http.StatusAccepted, response)

}

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	// mux.HandleFunc("/", makeHTTPHandlerFunc(logMiddleware(home)))
	mux.HandleFunc("POST /trip/preview", makeHTTPHandlerFunc(logMiddleware(handleTripPreview)))
	return mux
}




func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			log.Printf("handler error on %s %s: %v", r.Method, r.URL.Path, err)
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}