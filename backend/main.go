package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func streamResponse(w http.ResponseWriter, _ *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming response unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	for i := 1; i <= 50; i++ {
		fmt.Fprintf(w, "Data %d\n", i)
		flusher.Flush()
		time.Sleep(500 * time.Millisecond)
	}
}

func slowResponse(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Data from server at %s", time.Now())
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/fetch", slowResponse)
	r.Get("/stream", streamResponse)
	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("WORKING"))
	})

	fmt.Println("Server is running on port :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
