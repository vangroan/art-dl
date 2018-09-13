package artdl

import (
	"testing"
)

func TestPush(t *testing.T) {
	// Arrange
	q := NewQueue()

	// Act
	q.Push(0)
	q.Push(1)
	q.Push(2)

	// Assert
	expected := 3
	if q.Len() != expected {
		t.Fatalf("Expected %d, actual %d", q.Len(), expected)
	}
}

func TestPop(t *testing.T) {
	// Arrange
	q := NewQueue()
	data := []int{3, 5, 7, 9}
	for i := range data {
		q.Push(i)
	}

	for i := range data {
		// Act
		elem, ok := q.Pop()

		// Assert
		if !ok {
			t.Fatalf("Failed to pop from queue")
		}

		if elem != i {
			t.Fatalf("Expected %d, actual %d", i, elem)
		}
	}

	if q.Len() != 0 {
		t.Fatalf("Expected %d, actual %d", 0, q.Len())
	}
}

func TestPopEmpty(t *testing.T) {
	// Arrange
	q := NewQueue()

	// Act
	elem, ok := q.Pop()

	// Assert
	if elem != nil {
		t.FailNow()
	}

	if ok {
		t.FailNow()
	}
}

func TestPeekNotEmpty(t *testing.T) {
	// Arrange
	q := NewQueue()
	q.Push(1)

	// Act
	elem, ok := q.Peek()

	// Assert
	if elem == nil {
		t.Fail()
	}

	if !ok {
		t.Fail()
	}
}

func TestPeekEmpty(t *testing.T) {
	// Arrange
	q := NewQueue()

	// Act
	elem, ok := q.Peek()

	// Assert
	if elem != nil {
		t.Fail()
	}

	if ok {
		t.Fail()
	}
}
