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

// get query params
	if pathList, ok := query["p"]; ok { // path
		path := pathList[0]
		if path != "" {
			log.Printf("got path from query: %v\n", path)
		}
	}

// send response
	json.NewEncoder(writer).Encode("foo")
}

// $ curl "http://localhost:1000/v0/?p=/path/to/dir"
func findDupersHandler(writer http.ResponseWriter, request *http.Request){
	query := request.URL.Query()
	log.Printf("%v\n", query)

	findConfig := finder.FindConfig{}
	hashConfig := finder.HashConfig{}
	// formatConfig := finder.FormatConfig{}

	if pathList, ok := query["p"]; ok { // path
		path := pathList[0]
		if path != "" {
			log.Printf("got path from query: %v\n", path)

			dupes := finder.FindDupes(path, findConfig, hashConfig)
			// for _, entries := range dupes {
			// 	format := finder.DupesFormatter(entries, formatConfig)
			// 	fmt.Printf("%s", format) // format has newline embedded at the end
			// }
			json.NewEncoder(writer).Encode(dupes)
		}
	}


}

func handleRequests(){
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/v0/", findDupersHandler)
	log.Fatal(http.ListenAndServe(":1000", nil))
}

func main(){
	log.Println("hellooo world")
	handleRequests()
}
