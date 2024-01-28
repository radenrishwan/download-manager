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

	url := "https://cdn.discordapp.com/attachments/566486385377280001/1186878617951678464/Velocity-ArcaneMoon30.ppf?ex=659e1459&is=658b9f59&hm=10cf0ea6136a3fb015e17d3a2629a15809bb106c5d3fcfa6f14394f013730b28&"
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
