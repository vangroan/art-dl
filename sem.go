package main

import "fmt"

// Semaphores keeps a semaphore channel per domain
type Semaphores struct {
	inner map[string]chan struct{}
}

// NewSemaphores creates a new instance of a semaphore mapping.
func NewSemaphores() Semaphores {
	return Semaphores{
		inner: make(map[string]chan struct{}),
	}
}

// Add creates a semaphore associated with the key.
//
// Returns error when key already exists.
func (sem *Semaphores) Add(key string, concurrency int) error {
	if _, ok := sem.inner[key]; ok {
		return fmt.Errorf("Key %s already exists", key)
	} else {
		sem.inner[key] = make(chan struct{})
	}

	return nil
}
