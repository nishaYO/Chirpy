package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)

	cfg := apiConfig{
		fileserverHits: 0,
	}
	mux.Handle("/app/*", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", checkHealth)
	mux.HandleFunc("GET /admin/metrics", cfg.logMetrics)
	mux.HandleFunc("GET /api/reset", cfg.resetMetrics)
	mux.HandleFunc("POST /api/validate_chirp", validate_chirp_len)

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
	w.Header().Set("Content-Type", "text/html")
	metricsPage := `
	<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>`
	metrics := fmt.Sprintf(metricsPage, cfg.fileserverHits)
	w.Write([]byte(metrics))
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.Write([]byte(strconv.Itoa(cfg.fileserverHits)))
}

func validate_chirp_len(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type chirps struct {
		Body string `json:"body"`
	}

	type resErr struct {
		Error string `json:"error"`
	}

	type resValid struct {
		Valid bool `json:"valid"`
	}
	// decode request
	decoder := json.NewDecoder(r.Body)
	chirp := chirps{}
	err := decoder.Decode(&chirp)

	if err != nil {
		w.WriteHeader(500)
		res := resErr{Error: "Something went wrong"}
		dat, _ := json.Marshal(res)
		w.Write([]byte(dat))
		return
	}
	// check length of the chirp
	if len(chirp.Body) > 140 {
		w.WriteHeader(400)
		res := resErr{Error: "Chirp is too long"}
		dat, _ := json.Marshal(res)
		w.Write([]byte(dat))
		return
	}
	w.WriteHeader(200)
	res := resValid{Valid: true}
	dat, _ := json.Marshal(res)
	w.Write([]byte(dat))
}
