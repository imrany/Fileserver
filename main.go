package main

import (
	"net/http"
	"os"
	"log"
	"encoding/json"
)

type Message struct{
	Message string `json:"message"`
}

func helloWorldJson(w http.ResponseWriter , r *http.Request){
	helloMsg :=Message{
		Message:"Hello world",
	}

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(helloMsg)
}

func main(){
	fs := http.FileServer(http.Dir("./views"))

	router := http.NewServeMux()
	router.Handle("GET /views/", http.StripPrefix("/views/", fs))
	router.HandleFunc("GET /", helloWorldJson)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := http.Server{
		Addr: "0.0.0.0:" + port,
		Handler: router,
	}

	log.Printf("Server running on PORT %v", port)
	srv.ListenAndServe()
}
