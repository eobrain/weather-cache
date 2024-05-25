package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

var apiKey string

func init() {
	apiKey = os.Getenv("OPENWEATHERMAP_API_KEY")
	if apiKey == "" {
		log.Fatal("environment variable OPENWEATHERMAP_API_KEY is not set")
	}
}

func api(lat float64, lon float64) (string, error) {
	requestURL := fmt.Sprintf(`https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&units=metric&appid=%s`, lat, lon, apiKey)
	res, err := http.Get(requestURL)
	if err != nil {
		return fmt.Sprintf("error making http request: '%s'", requestURL), err
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	return string(resBody), nil
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

	w.Header().Set("Content-Type", "application/json")
	result, err := api(lat, lon)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(result))
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
