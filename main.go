package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func api(lat float64, lon float64) string {
	return fmt.Sprintf("{\"lat\": %f, \"lon\": %f}", lat, lon)
}

func Api(w http.ResponseWriter, req *http.Request) {

	query := req.URL.Query()
	lat, err := strconv.ParseFloat(query.Get("lat"), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	lon, err := strconv.ParseFloat(query.Get("lon"), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "text/jsom")
	w.Write([]byte(api(lat, lon)))
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
