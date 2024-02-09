package main

import "os"

func main() {
	Exit(1)
}

// Exit - тест вызова метода os.Exit из экспортируемой функции внутри метода main
func Exit(num int) {
	os.Exit(num)
}
