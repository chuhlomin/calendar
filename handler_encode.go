package main

import (
	"encoding/base64"
	"log"
	"net/http"

	"google.golang.org/protobuf/proto"
)

func handlerEncode(w http.ResponseWriter, r *http.Request) {
	in, err := parsePDFRequest(r)
	if err != nil {
		log.Printf("error parsing PDF request: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// marshal to protobuf
	b, err := proto.Marshal(&in.request)
	if err != nil {
		log.Printf("error marshaling proto: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Sorry, something went wrong"))
		return
	}

	// base64 encode the bytes
	b64 := base64.StdEncoding.EncodeToString(b)

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(b64))
}
