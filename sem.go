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
		sem.inner[key] = make(chan struct{}, concurrency)
	}

	return nil
}

// Acquire takes a slot from the semaphore associated
// with the key, blocking the current goroutine if there
// is no available slot.
//
// Returns an error if the key does not exist.
func (sem *Semaphores) Acquire(key string) error {
	if semaphore, ok := sem.inner[key]; ok {
		semaphore <- struct{}{}

		return nil
	} else {
		return fmt.Errorf("Key %s does not exist", key)
	}
}

// Release frees up a slot at the given semaphore, allowing
// other goroutines to resume execution.
//
// Returns an error if the key does not exist.
func (sem *Semaphores) Release(key string) error {
	if semaphore, ok := sem.inner[key]; ok {
		<-semaphore

		return nil
	} else {
		return fmt.Errorf("Key %s does not exist", key)
	}
}
