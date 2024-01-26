package main

import (
	"log"
	"net/http"
	"sync"

	downloadmanager "github.com/radenrishwan/download-manager"
)

func main() {
	var wg sync.WaitGroup
	mux := http.NewServeMux()

	mux.HandleFunc("/download", DownloadHandler)

	wg.Add(1)
	go func() {
		log.Println("Server started at localhost:3000")
		if err := http.ListenAndServe(":3000", mux); err != nil {
			wg.Done()
		}
	}()

	url := "http://localhost:3000/download"
	log.Println("Downloading file from", url)
	result, err := downloadmanager.GetMetaData(url)
	if err != nil {
		log.Fatalln(err)
	}

	err = downloadmanager.DownloadFile(result, 5)
	if err != nil {
		log.Fatalln(err)
	}

	wg.Wait()
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	filename := "demo.gif"
	mimetype := "text/plain"

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", mimetype)

	http.ServeFile(w, r, "dummy/"+filename)
}
