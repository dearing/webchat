// a simple webchat over websockets
package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"log"
	"net/http"
)

var config Config

var conf = flag.String("conf", "webchat.conf", "json configuration")
var gen = flag.Bool("gen", false, "generate a default conf in working directory")

func main() {

	flag.Parse()

	if *gen {
		config.GenerateConfig(*conf)
		return
	}

	config.LoadConfig(*conf)

	// spin up the router
	go Router.HandleClients()

	// assign our static fileserver and websocket to a mux
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(config.WWWRoot)))
	mux.Handle(config.Tag, websocket.Handler(WebsocketHandler))

	// start listening
	if config.Verbose {
		log.Printf("listening on %s // root=%s\n", config.WWWHost, config.WWWRoot)
	}

	if config.TLS {
		if err := http.ListenAndServeTLS(config.WWWHost, config.Certificate, config.CertificateKey, mux); err != nil {
			log.Println(err)
		}
	} else {
		if err := http.ListenAndServe(config.WWWHost, mux); err != nil {
			log.Println(err)
		}

	}

}
