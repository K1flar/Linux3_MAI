package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"task1/cprinter"
)

type Request struct {
	XMLName  xml.Name `xml:"file"`
	FilePath string   `xml:",chardata"`
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("no port")
		os.Exit(1)
	}
	port := os.Args[1]

	// init xml-rpc server
	http.HandleFunc("/rpc/print", print)

	if err := http.ListenAndServe(fmt.Sprintf("localhost:%s", port), nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func print(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	req := Request{}
	err = xml.Unmarshal(b, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = cprinter.Print(req.FilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
