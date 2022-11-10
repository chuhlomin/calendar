package main

import "net/http"

func handlerStatic(dir, file string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, dir+"/"+file)
	}
}
