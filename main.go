package main

import (
	"fmt"
	"sync"

	"github.com/labstack/echo/v5"
)

type Store[K comparable, V any] interface {
	Put(K, V) error
	Get(K) (V, error)
	Update(K, V) error
	Delete(K) error
}

type KVStore[K comparable, V any] struct {
	data map[K]V
	mu   sync.RWMutex
}

// Constructor
func InitKVStore[K comparable, V any]() KVStore[K, V] {
	return KVStore[K, V]{
		data: make(map[K]V),
	}
}

// Checks if given key exists
func (s *KVStore[K, V]) Exists(key K) bool {
	_, exists := s.data[key]
	return exists
}

func (s *KVStore[K, V]) Put(key K, value V) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Exists(key) {
		return fmt.Errorf("ERROR: The Key (%v) already exists!", key)
	}
	s.data[key] = value
	return nil
}

func (s *KVStore[K, V]) Get(key K) (V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, exists := s.data[key]
	if !exists {
		return val, fmt.Errorf("ERROR: The Key (%v) does not exist!", key)
	}

	return val, nil
}

func (s *KVStore[K, V]) Update(key K, value V) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Exists(key) {
		return fmt.Errorf("ERROR: The Key (%v) does not exist!", key)
	}

	s.data[key] = value
	return nil
}

func (s *KVStore[K, V]) Delete(key K) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Exists(key) {
		return fmt.Errorf("ERROR: The Key (%v) does not exist!", key)
	}

	delete(s.data, key)
	return nil
}

// Server
type Server struct {
	Store      KVStore[string, string]
	ListenAddr string
}

// Creates a new server instance
func NewServer(listenAddr string) *Server {
	return &Server{
		Store:      InitKVStore[string, string](),
		ListenAddr: listenAddr,
	}
}

// Starts the HTTP server
func (s *Server) Start() {
	fmt.Printf("HTTP server is running on port %s", s.ListenAddr)
	e := echo.New()

	e.GET("/put/:key/:value", s.HandlePut)
	e.GET("/get/:key", s.HandleGet)
	e.GET("/update/:key/:value", s.HandleUpdate)
	e.GET("/delete/:key", s.HandleDelete)

	e.Start(s.ListenAddr)
}

func main() {
	s := NewServer(":3000")
	s.Start()
}
