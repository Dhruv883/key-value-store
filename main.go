package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/labstack/echo/v5"
)

type Store[K comparable, V any] interface {
	Put(K, V) error
	PutWithTTL(K, V, time.Duration) error
	Get(K) (V, error)
	Update(K, V) error
	Delete(K) error
	SetTTL(K, time.Duration) error
	TTLRemaining(K) (time.Duration, bool, error)
}

type Entry[V any] struct {
	value     V
	expiresAt time.Time
	hasTTL    bool
}

func (e Entry[V]) isExpired() bool {
	return e.hasTTL && time.Now().After(e.expiresAt)
}

type KVStore[K comparable, V any] struct {
	data map[K]Entry[V]
	mu   sync.RWMutex
}

// Constructor
func InitKVStore[K comparable, V any]() KVStore[K, V] {
	return KVStore[K, V]{
		data: make(map[K]Entry[V]),
	}
}

// Checks if given key exists
func (s *KVStore[K, V]) Exists(key K) bool {
	entry, exists := s.data[key]
	if !exists {
		return false
	}
	if entry.isExpired() {
		return false
	}
	return true
}

func (s *KVStore[K, V]) Put(key K, value V) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Exists(key) {
		return fmt.Errorf("ERROR: The Key (%v) already exists!", key)
	}
	s.data[key] = Entry[V]{value: value}
	return nil
}

func (s *KVStore[K, V]) PutWithTTL(key K, value V, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Exists(key) {
		return fmt.Errorf("ERROR: The Key (%v) already exists!", key)
	}
	s.data[key] = Entry[V]{value: value, expiresAt: time.Now().Add(ttl), hasTTL: true}
	return nil
}

func (s *KVStore[K, V]) Get(key K) (V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.data[key]
	if !exists || entry.isExpired() {
		var zero V
		return zero, fmt.Errorf("ERROR: The Key (%v) does not exist!", key)
	}

	return entry.value, nil
}

func (s *KVStore[K, V]) Update(key K, value V) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Exists(key) {
		return fmt.Errorf("ERROR: The Key (%v) does not exist!", key)
	}

	e := s.data[key]
	e.value = value
	s.data[key] = e

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

// Returns the remaining TTL for a key, a bool indicating whether the key has a TTL, and any error.
func (s *KVStore[K, V]) TTLRemaining(key K) (time.Duration, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.data[key]
	if !exists || entry.isExpired() {
		return 0, false, fmt.Errorf("ERROR: The Key (%v) does not exist!", key)
	}

	if !entry.hasTTL {
		return 0, false, nil
	}

	return time.Until(entry.expiresAt), true, nil
}

func (s *KVStore[K, V]) SetTTL(key K, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Exists(key) {
		return fmt.Errorf("ERROR: The Key (%v) does not exist!", key)
	}

	e := s.data[key]
	e.expiresAt = time.Now().Add(ttl)
	e.hasTTL = true
	s.data[key] = e
	return nil
}

func (s *KVStore[K, V]) Cleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			s.mu.Lock()
			for key, e := range s.data {
				if e.isExpired() {
					delete(s.data, key)
				}
			}
			s.mu.Unlock()
		}
	}()
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
	s.Store.Cleanup(time.Minute)

	e := echo.New()

	e.GET("/put/:key/:value", s.HandlePut)
	e.GET("/get/:key", s.HandleGet)
	e.GET("/update/:key/:value", s.HandleUpdate)
	e.GET("/delete/:key", s.HandleDelete)
	e.GET("/ttl/:key/:seconds", s.HandleSetTTL)
	e.GET("/ttl/:key", s.HandleGetTTL)

	e.Start(s.ListenAddr)
}

func main() {
	s := NewServer(":3000")
	s.Start()
}
