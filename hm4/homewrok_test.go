package main

import (
	"golang.org/x/exp/constraints"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Идея упорядоченного словаря заключается в том, что он будет реализован на основе бинарного дерева поиска
Дерево будет строиться только по ключам элементов, значения элементов при построении дерева не учитываются.
Элементы с одинаковыми ключами в упорядоченном словаре хранить нельзя.
Поподробнее с бинарными деревьями поиска можно познакомиться https://habr.com/ru/articles/65617/

type OrderedMap struct { ... }

func NewOrderedMap() OrderedMap                      // создать упорядоченный словарь
func (m \*OrderedMap) Insert(key, value int)          // добавить элемент в словарь
func (m \*OrderedMap) Erase(key int)                  // удалить элемент из словари
func (m \*OrderedMap) Contains(key int) bool          // проверить существование элемента в словаре
func (m \*OrderedMap) Size() int                      // получить количество элементов в словаре
func (m \*OrderedMap) ForEach(action func(int, int))  // применить функцию к каждому элементу словаря от меньшего к большему

*/
// go test -v homework_test.go

type OrderedMap[k constraints.Ordered, v any] struct {
	root *node[k, v]
	size int
}

type node[k constraints.Ordered, v any] struct {
	left, right *node[k, v]
	key         k
	value       v
}

func NewOrderedMap[k constraints.Ordered, v any]() *OrderedMap[k, v] {
	return &OrderedMap[k, v]{}
}

func (m *OrderedMap[k, v]) Insert(key k, value v) {
	if m.root == nil {
		m.root = &node[k, v]{key: key, value: value}
		m.size++

		return
	} else {
		m.size += m.insert(m.root, key, value)
	}
}

func (m *OrderedMap[k, v]) insert(n *node[k, v], key k, val v) int {
	if key == n.key {
		n.value = val

		return 0
	}
	if key < n.key {
		if n.left == nil {
			n.left = &node[k, v]{key: key, value: val}

			return 1
		}
		return m.insert(n.left, key, val)
	} else {
		if n.right == nil {
			n.right = &node[k, v]{key: key, value: val}

			return 1
		}

		return m.insert(n.right, key, val)
	}
}

func (m *OrderedMap[k, v]) Erase(key k) {
	m.root = m.erase(m.root, key)

}

func (m *OrderedMap[k, v]) erase(n *node[k, v], key k) *node[k, v] {
	if n == nil {

		return nil
	}
	if key < n.key {
		n.left = m.erase(n.left, key)

		return n
	} else if key > n.key {
		n.right = m.erase(n.right, key)

		return n
	} else {
		if n.left == nil {
			m.size--

			return n.right
		}
		if n.right == nil {
			m.size--

			return n.left
		}
		minNode := m.minNode(n.right)
		n.key = minNode.key
		n.value = minNode.value
		n.right = m.erase(n.right, minNode.key)
		m.size--

		return n
	}
}

func (m *OrderedMap[k, v]) minNode(n *node[k, v]) *node[k, v] {
	if n.left == nil {
		return n
	}

	return m.minNode(n.left)
}

func (m *OrderedMap[k, v]) Contains(key k) bool {
	cur := m.root
	for cur != nil {
		if key < cur.key {
			cur = cur.left
		} else if key > cur.key {
			cur = cur.right
		} else {

			return true
		}
	}

	return false
}

func (m *OrderedMap[k, v]) Size() int {
	return m.size
}

func (m *OrderedMap[k, v]) ForEach(action func(k, v)) {
	m.forEach(m.root, action)
}

func (m *OrderedMap[k, v]) forEach(n *node[k, v], action func(k, v)) {
	if n == nil {
		return
	}
	m.forEach(n.left, action)
	action(n.key, n.value)
	m.forEach(n.right, action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap[int, int]()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
