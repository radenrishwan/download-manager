package downloadmanager

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// [DownloadMetaData] is a struct that contains metadata of a url to be downloaded
type DownloadMetaData struct {
	Url string
	// Length is the size of the file in bytes
	Length uint64
	// FileName is the name of the file
	FileName string
	// ContentType is the type of the file
	ContentType string
}

// [GetMetaData] is a function to get metadata of a url to be downloaded
func GetMetaData(url string) (DownloadMetaData, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return DownloadMetaData{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return DownloadMetaData{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode > 300 {
		return DownloadMetaData{}, errors.New("could not get metadata, maybe the URL is wrong?")
	}

	contentLength, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return DownloadMetaData{}, errors.New("could not get metadata, maybe the URL is wrong?")
	}

	filename := resp.Header.Get("Content-Disposition")
	contentType := resp.Header.Get("Content-Type")

	if filename == "" {
		filename = time.Now().String() + "." + contentType
	}

	return DownloadMetaData{
		Url:         url,
		Length:      uint64(contentLength),
		FileName:    filename,
		ContentType: contentType,
	}, nil
}

type rangeHeader struct {
	Start uint64
	End   uint64
}

// [DownloadFile] is a function to download a file from a url
func DownloadFile(metadata DownloadMetaData, parrarel int) error {
	var header []rangeHeader
	var currentPart uint64 = 0

	for i := 0; i < parrarel; i++ {
		fileSize := metadata.Length / uint64(parrarel)
		header = append(header, rangeHeader{
			Start: currentPart,
			End:   currentPart + fileSize - 1,
		})

		if i == parrarel-1 {
			fileSize += metadata.Length % uint64(parrarel)

			header[i].End += metadata.Length % uint64(parrarel)
		}

		currentPart += fileSize
	}

	filename := make(chan string, parrarel)
	var res []string
	for i := 0; i < parrarel; i++ {
		go download(filename, metadata, header[i], i)
	}

	// get all the file names from the channel
	for v := range filename {
		res = append(res, v)
		if len(res) == parrarel {
			sort.Strings(res)
			break
		}
	}

	f, err := os.OpenFile(strings.Split(metadata.FileName, "=")[1], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	for _, v := range res {
		log.Println("merging file", v)
		file, err := os.Open(v)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		file.Close()
		os.Remove(v)
	}

	return nil
}

type writeCounter struct {
	Filename string
	Total    uint64
}

// Write implements the io.Writer interface.
//
// Always completes and never returns an error.
func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	fmt.Printf("Read %d bytes for a total of %d for %s\n", n, wc.Total, wc.Filename)
	return n, nil
}

func download(result chan<- string, metadata DownloadMetaData, header rangeHeader, part int) error {
	req, err := http.NewRequest("GET", metadata.Url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Range", "bytes="+strconv.Itoa(int(header.Start))+"-"+strconv.Itoa(int(header.End)))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode > 299 {
		return errors.New("could not download file, maybe the URL is wrong?")

	}

	filename := strings.Split(metadata.FileName, "=")[1] + ".part" + strconv.Itoa(part) + ".🔥"

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	src := io.TeeReader(resp.Body, &writeCounter{
		Filename: filename,
	})

	// write file
	_, err = io.Copy(file, src)
	if err != nil {
		return err
	}

	result <- filename

	defer file.Close()
	defer resp.Body.Close()
	return nil
}
