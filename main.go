package main

import (
	"fmt"
	"log"
	"net/http"

	"ascii-art-web-export/source"
)

const PORT = ":1010"

func main() {
	http.HandleFunc("/css/", source.StaticHandler)
	http.HandleFunc("/", source.MainHandler)
	http.HandleFunc("/ascii-art-web", source.HandleAscii)
	http.HandleFunc("/download/", source.Download)
	http.HandleFunc("/downloadfile", source.DownloadFile)
	// http.HandleFunc("/ascii-art-web", source.HandleAscii)
	// http.HandleFunc("/ascii-art-web", source.HandleAscii)
	
	fmt.Println("http://localhost" + PORT)
	
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
