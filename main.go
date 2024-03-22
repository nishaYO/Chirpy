package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)
	mux.Handle("/", http.FileServer(http.Dir(".")))
	server := http.Server{
		Addr:    ":8000",
		Handler: corsMux,
	}
	server.ListenAndServe()
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
