package main

import (
	"log"
	"net/http"
	"sync"
)

type Stream struct {
	mu      sync.RWMutex
	status  string
	headers http.Header
	body    []byte
	clients map[chan []byte]struct{}
	done    chan struct{}
}

type Manager struct {
	mu      sync.RWMutex
	streams map[string]*Stream
}

func NewManager() *Manager {
	return &Manager{
		streams: make(map[string]*Stream),
	}
}

func (m *Manager) GetOrCreateStream(key string) (*Stream, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	stream, exists := m.streams[key]
	if !exists {
		stream = &Stream{
			status:  "init",
			headers: make(http.Header),
			clients: make(map[chan []byte]struct{}),
			done:    make(chan struct{}),
		}
		m.streams[key] = stream
	}

	return stream, exists
}

func (m *Manager) RemoveStream(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.streams, key)
}

func (s *Stream) AddClient(client chan []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients[client] = struct{}{}
}

func (s *Stream) RemoveClient(client chan []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.clients, client)
}

func (s *Stream) Broadcast(data []byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for client := range s.clients {
		client <- data
	}
}

func handleRequest(m *Manager, key string, stream *Stream) {
	defer m.RemoveStream(key)
	defer close(stream.done)
}

func main() {
	manager := NewManager()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path
		stream, exists := manager.GetOrCreateStream(key)

		if exists {
			go handleRequest(manager, key, stream)
		}

		client := make(chan []byte, 100)
		stream.AddClient(client)
		defer stream.RemoveClient(client)

		for k, v := range stream.headers {
			w.Header()[k] = v
		}

		for {
			select {
			case data := <-client:
				_, err := w.Write(data)
				if err != nil {
					log.Printf("error writing to client: %v", err)
					return
				}
				w.(http.Flusher).Flush()
			case <-stream.done:
				return
			case <-r.Context().Done():
				return
			}
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
