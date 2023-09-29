package main

import (
	"errors"
	"flag"
	"fmt"
	"image/png"
	"io"
	"net/http"
	"os"
)

func main() {

	contentPathPrefix := flag.String("prefix", "/Users/i.nedzvetskiy/go/src/dns-image-gen", "content")
	origin := flag.String("host", "127.0.0.1:3333", "[host]:port")
	flag.Parse()

	g := NewGrapher(*contentPathPrefix)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = io.WriteString(writer, "Hello world!\n")
	})
	mux.HandleFunc("/img", func(writer http.ResponseWriter, request *http.Request) {
		domain := request.URL.Query().Get("d")
		dc := g.Draw(domain)
		err := png.Encode(writer, dc.Image())
		if err != nil {
			fmt.Printf("Error on send image %v\n", err)
		}
	})
	//dc := g.Draw("gems")
	//
	//dc.SavePNG("/Users/i.nedzvetskiy/go/src/dns-image-gen/out.png")

	fmt.Printf("Starting %s\n", *origin)
	err := http.ListenAndServe(*origin, mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
