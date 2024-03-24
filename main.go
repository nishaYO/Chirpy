package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)

	cfg := &apiConfig{} //create instance of api config struct
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("/healthz", checkHealth)
	mux.HandleFunc("/metrics", cfg.logMetrics)
	mux.HandleFunc("/reset", cfg.resetMetrics)

	server := http.Server{
		Addr:    ":8080",
		Handler: corsMux,
	}
	server.ListenAndServe()
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) logMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := fmt.Sprintf("Hits: %d", cfg.fileserverHits)
	w.Write([]byte(metrics))
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.Write([]byte(strconv.Itoa(cfg.fileserverHits)))
}
