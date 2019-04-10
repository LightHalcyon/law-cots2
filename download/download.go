package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/reznov53/law-cots2/mq"
)

var ch *mq.Channel
var rKey string
var err error

// Source: https://golangcode.com/download-a-file-with-progress/

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer
// interface and we can pass this into io.TeeReader() which will report progress on each
// write cycle.
type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.printProgress()
	return n, nil
}

func (wc WriteCounter) printProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	s := fmt.Sprintf("\r%s", strings.Repeat(" ", 35))

	ch.PostMessage(s, rKey)

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	s = fmt.Sprintf("\rDownloading... %s complete", humanize.Bytes(wc.Total))

	ch.PostMessage(s, rKey)
}

// File will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory. We pass an io.TeeReader
// into Copy() to report progress on the download.
func File(filepath string, url string, channel *mq.Channel, routeKey string) error {

	rKey = routeKey
	ch = channel

	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return err
	}

	// The progress use the same line so print a new line once it's finished downloading
	// fmt.Print("\n")

	err = os.Rename(filepath+".tmp", filepath)
	if err != nil {
		return err
	}

	ch.PostMessage("Download Finished", rKey)

	return nil
}

// func init() {
// 	// url := "amqp://" + os.Getenv("UNAME") + ":" + os.Getenv("PW") + "@" + os.Getenv("URL") + ":" + os.Getenv("PORT") + "/"
// 	url = "amqp://1406568753:167664@152.118.148.103:5672/"
// 	// vhost := os.Getenv("VHOST")
// 	vhost = "1406568753"
// 	// exchangeName := os.Getenv("EXCNAME")
// 	exchangeName = "1406568753"
// 	exchangeType = "direct"
// 	ch, err = InitMQ(url, vhost, exchangeName, exchangeType)
// 	if err != nil {
// 		panic(err)
// 	}
// }