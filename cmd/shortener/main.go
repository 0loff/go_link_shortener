package main

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
)

var linkStorage = map[string]string{}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", requestHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createLink(w, r)

	case http.MethodGet:
		getLink(w, r)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func createLink(w http.ResponseWriter, r *http.Request) {
	rData, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	if len(r.URL.Path[1:]) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	linkStorage[base64.RawStdEncoding.EncodeToString(rData)] = string(rData)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://" + r.Host + "/" + base64.RawStdEncoding.EncodeToString(rData)))
}

func getLink(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) > 1 {
		link, ok := linkStorage[r.URL.Path[1:]]

		if ok {
			w.Header().Set("Location", link)
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
