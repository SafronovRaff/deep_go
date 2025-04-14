package hm1

import (
	"unsafe"
)

// **Нужно реализовать функцию по конвертации числа из прямого порядка следования байт (*Big Endian*) в обратный порядок следования байт (*Little Endian*).**
//
//
// Например, число `0x01020304` при вызове функции должно быть сконвертировано в число `0x04030201`, а число `0x0000FFFF` в `0xFFFF0000`.
// Поподробнее про порядки следования байт можно прочитать [здесь](https://betterexplained.com/articles/understanding-big-and-little-endian-byte-order/).

// go test -v homework_test.go

// Эффективная работа с битами при помощи Go https://habr.com/ru/companies/ruvds/articles/744230/
// Хитрости с битовыми операциями https://tproger.ru/articles/awesome-bits

func ToLittleEndian(number uint32) (res uint32) {
	b1 := (number >> 24) & 0xFF
	b2 := (number >> 16) & 0xFF
	b3 := (number >> 8) & 0xFF
	b4 := number & 0xFF

	res = (b4 << 24) | (b3 << 16) | (b2 << 8) | b1

	return res
}

func ToLittleEndian2(number uint32) (res uint32) {
	const octetsCount = 4 // (32 бита / 8 бит на байт и того 4 октета).
	p := unsafe.Pointer(&number)

	for idx := 0; idx < octetsCount; idx++ {

		byteValue := *(*uint8)(unsafe.Add(p, idx))

		bitOffset := (octetsCount - idx - 1) * 8
		res |= uint32(byteValue) << bitOffset
	}

	return res
}

type Number interface {
	uint16 | uint32 | uint64
}

func ToLittleEndian3[T Number](number T) T {
	var result T
	octetsCount := unsafe.Sizeof(number)
	p := unsafe.Pointer(&number)

	for idx := uintptr(0); idx < octetsCount; idx++ {
		byteValue := *(*uint8)(unsafe.Add(p, idx))
		bitOffset := (octetsCount - idx - 1) * 8
		result |= T(byteValue) << bitOffset
	}

	return result
}
