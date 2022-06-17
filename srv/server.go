package main

import (
	"log"
	"net/http"
	"encoding/json"
	"dupefinder/src" // "dupefinder/src" as finder
	// "fmt"
)

// $ curl "http://localhost:1000/?p=/path/to/dir"
func rootHandler(writer http.ResponseWriter, request *http.Request){
	query := request.URL.Query()
	log.Printf("rootHandler got query: %v\n", query)

// send response
	json.NewEncoder(writer).Encode("foo")
}

// $ curl "http://localhost:1000/v0/?p=/path/to/dir"
func findDupersHandler(writer http.ResponseWriter, request *http.Request){
	query := request.URL.Query()
	log.Printf("%v\n", query)

	findConfig := finder.FindConfig{}
	hashConfig := finder.HashConfig{}

	if pathList, ok := query["p"]; ok { // path
		path := pathList[0]
		if path != "" { // map[p:[/path/to/dir]] ; default to empty string
			log.Printf("got path from query: %v\n", path)
			dupes := finder.FindDupes(path, findConfig, hashConfig)
			json.NewEncoder(writer).Encode(dupes)
		}
	}


}

func handleRequests(port string){
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/v0/", findDupersHandler)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}

func main(){
	port := "1000"
	url := "http://localhost"
	log.Printf("Running on %v:%v\n", url, port)
	handleRequests(port)
}
