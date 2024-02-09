package main

import "os"

// Interrupter - тестовая структура используемая в качестве ресивера
// для тестового вызова метода os.Exit из метода main
type Interrupter struct {
}

// Exit - метод структуры Interrupter для тестирования вызова os.Exit из метода main
func (i Interrupter) Exit(code int) {
	os.Exit(code)
}

func main() {
	t := Interrupter{}
	t.Exit(1)
}
