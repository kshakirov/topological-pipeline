//go:build !debug

package main

// Пустая функция — компилятор Go её полностью вырезает (inline)
func Debug(msg string, args ...any) {}
