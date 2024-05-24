package main

import (
	// "fmt"
	// "io"
	"log"
	"net/http"
)

func Api(w http.ResponseWriter, req *http.Request) {

	query := req.URL.Query()
	lat := query.Get("lat")
	lon := query.Get("lon")

	w.Header().Set("Content-Type", "text/jsom")
	w.Write([]byte("{\"lat\": " + lat + ", \"lon\": " + lon + "}"))
	// fmt.Fprintf(w, "This is an example server.\n")
	// io.WriteString(w, "This is an example server.\n")
}

func main() {
	http.HandleFunc("/api", Api)
	err := http.ListenAndServeTLS(":1443", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
