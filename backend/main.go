package main

import (
	"fmt"
	"net/http"
	"time"
)

func streamResponse(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming response unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	for i := 1; i <= 10; i++ {
		fmt.Fprintf(w, "Data %d\n", i)
		flusher.Flush()
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	http.HandleFunc("/", streamResponse)
	fmt.Println("Server is running on port :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
