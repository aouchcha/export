package main

import (
	"fmt"
	"net/http"

	f "ascii-art-web-export/func"
)

func main() {
	http.HandleFunc("/styles/", f.ServeStyle)
	http.HandleFunc("/", f.Welcom)
	//http.HandleFunc("/check",f.Check)
	http.HandleFunc("/ascii-art", f.Last)
	//http.HandleFunc("/download",f.Download)
	http.HandleFunc("/output", f.Output)
	fmt.Println("the server is running on localhost port 9099")
	fmt.Println("http://localhost:9099")
	err := http.ListenAndServe(":9099", nil)
	if err != nil {
		fmt.Println(err)
	}
}
