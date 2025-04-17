package main

import (
	"reflect"
	"slices"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// Идея подхода copy-on-write заключается в том, что при чтении данных используется общая копия данных буффера,
// но в случае изменения данных — создается новая копия данных буффера.
// Для реализации такого подхода можно использовать разделяемый счетчик ссылок - если при изменении данных буффера кто-то еще ссылается на этот буффер,
// то нужно будет сначала произвести копию данных буффера, изменить счетчик ссылок и только затем произвести изменение
// (*если никто не ссылается на буффер, то копировать данные буффера не нужно при изменении данных*).

// Дополнительно еще нужно реализовать метод конвертации данных буффера в строку без копирования и дополнительного выделения памяти.
// Предположим, что все будут производить копирование
// буффера только с использованием метода Clone()
// https://habr.com/ru/articles/673372/
// https://www.lenovo.com/us/en/glossary/what-is-cow/?orgRef=https%253A%252F%252Fwww.google.com%252F&srsltid=AfmBOoo8DwtmdYLBNkIA62JQ82t9bFSPMg26sPfuVPLPrnmbY9BnF4xy
// https://gayeugur.medium.com/what-is-copy-on-write-cow-in-swift-b51f269cfe21

type COWBuffer struct {
	data []byte
	refs *int
}

// создать буффер с определенными данными
func NewCOWBuffer(data []byte) COWBuffer {
	cnt := new(int)
	*cnt = 1

	return COWBuffer{
		data: data,
		refs: cnt,
	}
}

// создать новую копию буфера
func (b *COWBuffer) Clone() COWBuffer {
	*b.refs++

	return *b
}

// перестать использовать копию буффера
func (b *COWBuffer) Close() {
	*b.refs--
}

// изменить определенный байт в буффере
func (b *COWBuffer) Update(index int, value byte) bool {
	if b == nil {
		return false
	}

	if index < 0 || index >= len(b.data) {
		return false
	}

	if *b.refs > 1 {
		*b.refs--
		*b = NewCOWBuffer(slices.Clone(b.data))
	}

	b.data[index] = value

	return true
}

// сконвертировать буффер в строку
func (b *COWBuffer) String() string {
	return unsafe.String(unsafe.SliceData(b.data), len(b.data))
}

func TestCOWBuffer(t *testing.T) {
	data := []byte{'a', 'b', 'c', 'd'}
	buffer := NewCOWBuffer(data)
	defer buffer.Close()

	copy1 := buffer.Clone()
	copy2 := buffer.Clone()

	assert.Equal(t, unsafe.SliceData(data), unsafe.SliceData(buffer.data))
	assert.Equal(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	assert.True(t, (*byte)(unsafe.SliceData(data)) == unsafe.StringData(buffer.String()))
	assert.True(t, (*byte)(unsafe.StringData(buffer.String())) == unsafe.StringData(copy1.String()))
	assert.True(t, (*byte)(unsafe.StringData(copy1.String())) == unsafe.StringData(copy2.String()))

	assert.True(t, buffer.Update(0, 'g'))
	assert.False(t, buffer.Update(-1, 'g'))
	assert.False(t, buffer.Update(4, 'g'))

	assert.True(t, reflect.DeepEqual([]byte{'g', 'b', 'c', 'd'}, buffer.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy1.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy2.data))

	assert.NotEqual(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	copy1.Close()

	previous := copy2.data
	copy2.Update(0, 'f')
	current := copy2.data

	// 1 reference - don't need to copy buffer during update
	assert.Equal(t, unsafe.SliceData(previous), unsafe.SliceData(current))

	copy2.Close()
}
