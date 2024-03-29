package adapter

import "net/http"

func ServeFile(f string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, f)
	})
}

func ServeFiles(p, d string) http.Handler {
	return http.StripPrefix(p, http.FileServer(http.Dir(d)))
}
