package artdl

import (
	"sync"
)

// Queue is a thread safe, FIFO queue.
//
// Internally it uses a slice, simply appending on push, and removing the first
// item when popping. It relies on slice behaviour where removing the first item
// does not reallocate the underlying array.
type Queue struct {
	data []interface{}
	lock *sync.RWMutex // pointer to avoid copy
}

// NewQueue creates a new empty queue.
func NewQueue() *Queue {
	return &Queue{
		data: make([]interface{}, 0),
		lock: &sync.RWMutex{},
	}
}

// Push adds an item to the end of the queue.
func (q *Queue) Push(item interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.data = append(q.data, item)
}

// Pop removes and returns the first item in the queue.
//
// If the queue is empty, the second return value will be false.
func (q *Queue) Pop() (interface{}, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.data) <= 0 {
		return nil, false
	}

	item := q.data[0]
	// Ensure the pointer is dropped so the contents can be collected
	q.data[0] = nil
	q.data = q.data[1:]
	return item, true
}

// Peek returns the last element in the queue.
//
// If the queue is empty, the second return value will be false.
func (q *Queue) Peek() (interface{}, bool) {
	q.lock.RLock()
	defer q.lock.RUnlock()

	if len(q.data) > 0 {
		return q.data[len(q.data)-1], true
	}
	return nil, false
}

// Len returns the number of elements in the queue.
func (q *Queue) Len() int {
	q.lock.RLock()
	defer q.lock.RUnlock()

	return len(q.data)
}
