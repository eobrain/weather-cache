package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	cache "github.com/patrickmn/go-cache"
)

var apiKey string
var c *cache.Cache

func init() {
	apiKey = os.Getenv("OPENWEATHERMAP_API_KEY")
	if apiKey == "" {
		log.Fatal("environment variable OPENWEATHERMAP_API_KEY is not set")
	}

	c = cache.New(24*time.Hour, 1*time.Hour)
}

func api(latScaled int, lonScaled int) (string, error) {
	cacheKey := fmt.Sprintf("%d:%d", latScaled, lonScaled)
	fmt.Printf("cacheKey=%s\n", cacheKey)

	result, found := c.Get(cacheKey)
	if found {
		fmt.Print("HIT\n")
		return result.(string), nil
	}
	fmt.Print("MISS\n")

	requestURL := fmt.Sprintf(
		`https://api.openweathermap.org/data/2.5/forecast?lat=%f&lon=%f&units=metric&appid=%s`,
		unscale(latScaled), unscale(lonScaled), apiKey)
	fmt.Printf("%s\n", requestURL)
	res, err := http.Get(requestURL)
	if err != nil {
		return fmt.Sprintf("error making http request: '%s'", requestURL), err
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "error reading response body", err
	}
	json := string(resBody)
	c.Set(cacheKey, json, cache.DefaultExpiration)
	return json, nil
}

func scale(x float64) int {
	return int(math.Round(x * 10000))
}
func unscale(i int) float64 {
	return float64(i) / 10000
}

func Api(w http.ResponseWriter, req *http.Request) {

	query := req.URL.Query()
	lat, err := strconv.ParseFloat(query.Get("lat"), 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	lon, err := strconv.ParseFloat(query.Get("lon"), 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	result, err := api(scale(lat), scale(lon))
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
