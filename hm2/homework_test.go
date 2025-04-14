package main

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

//**Кольцевая очередь** (*Circular Queue*) — это структура данных, которая представляет собой очередь (*FIFO*) фиксированного размера.
//Кольцевая очередь использует буфер фиксированного размера таким образом, как будто бы после последнего элемента
//сразу же снова идет первый (*как представлено на картинке справа*).
//Подробнее можно прочитать [здесь](https://www.programiz.com/dsa/circular-queue).
//
//
//Такая структура много где используется, например для организации различных очередей сообщений и
//а также буффер в буфферезированных каналах Go реализован в виде кольцевой очереди.
// https://evileg.com/ru/post/472/
// https://codechick.io/tutorials/dsa/dsa-circular-queue
// go test -v homework_test.go

type CircularQueue[T comparable] struct {
	values []T
	head   int
	tail   int
	size   int
	count  int
}

func NewCircularQueue[T comparable](size int) *CircularQueue[T] {
	return &CircularQueue[T]{
		values: make([]T, size),
		size:   size,
	}
}

func (q *CircularQueue[T]) Push(value T) bool {
	if q.Full() {
		return false
	}
	q.values[q.tail] = value
	q.tail = (q.tail + 1) % q.size
	q.count++

	return true
}

func (q *CircularQueue[T]) Pop() bool {
	if q.Empty() {
		return false
	}
	q.values[q.head] = zero[T]()
	q.head = (q.head + 1) % q.size
	q.count--

	return true
}

func (q *CircularQueue[T]) Front() (T, bool) {
	if q.Empty() {
		return zero[T](), false
	}

	return q.values[q.head], true
}

func (q *CircularQueue[T]) Back() (T, bool) {
	if q.Empty() {
		return zero[T](), false
	}
	back := (q.tail - 1 + q.size) % q.size

	return q.values[back], true
}

func (q *CircularQueue[T]) Empty() bool {
	return q.count == 0
}

func (q *CircularQueue[T]) Full() bool {
	return q.count == q.size
}

func zero[T any]() T {
	var z T
	return z
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3

	queue := NewCircularQueue[int](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	_, ok := queue.Front()
	require.False(t, ok)
	_, ok = queue.Back()
	require.False(t, ok)
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	front, ok := queue.Front()
	assert.Equal(t, 1, front)

	back, ok := queue.Back()
	assert.Equal(t, 3, back)

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	front, ok = queue.Front()
	assert.Equal(t, 2, front)

	back, ok = queue.Back()
	assert.Equal(t, 4, back)

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

}
