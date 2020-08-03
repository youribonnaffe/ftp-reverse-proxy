package main

import (
	"flag"
	"fmt"
	"github.com/jlaffaye/ftp"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {
	configuration := parseFlags()

	http.HandleFunc("/", proxy(configuration))

	log.Printf("Running on port %d, targetting %s ...", configuration.port, configuration.target.Host)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", configuration.port), nil))
}

type configuration struct {
	port   int
	target url.URL
}

func parseFlags() configuration {
	var target = flag.String("target", "", "FTP URL to proxy to (example: ftp://user:pwd@host:21)")
	var port = flag.Int("port", 8080, "Port to listen to")

	flag.Usage = func() {
		fmt.Printf("ftp-reverse-proxy, a HTTP reverse proxy to access a FTP server.\n" +
			"Usage:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *target == "" {
		flag.PrintDefaults()
		log.Fatal("Target is not defined")
	}

	var targetUrl, err = url.Parse(*target)
	if err != nil {
		flag.PrintDefaults()
		log.Fatal("Target URL is malformed")
	}

	return configuration{
		port:   *port,
		target: *targetUrl,
	}
}

func proxy(configuration configuration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := ftp.Dial(configuration.target.Host, ftp.DialWithTimeout(5*time.Second))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		password, _ := configuration.target.User.Password()
		err = c.Login(configuration.target.User.Username(), password)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if r.Method == "POST" {
			err = c.Stor(r.URL.Path, r.Body)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			log.Printf("Stored %s", r.URL.Path)

		} else if r.Method == "GET" {

			_, ok := r.URL.Query()["ls"]

			if ok {
				fileList, err := c.NameList(r.URL.Path)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusNotFound)
					return
				}
				for _, fileName := range fileList {
					fmt.Fprintf(w, fileName)
					fmt.Fprintf(w, "\n")
				}
				return
			}

			response, err := c.Retr(r.URL.Path)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusNotFound)
				return
			}

			_, err = io.Copy(w, response)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Printf("Retrieved %s", r.URL.Path)

		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}
}
