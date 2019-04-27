package demo

import (
	"testing"
)

func Fib(n int) int {
	if n < 2 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}

func FibSlow(n int) int {
	if n < 2 {
		return n
	}
	fibs := make([]int, 0)
	fibs = append(fibs, Fib(n-1), Fib(n-2))
	sum := 0
	for _, f := range fibs {
		sum += f
	}
	return sum
}

func BenchmarkFib10(b *testing.B) {
	b.Run("Fib()", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			Fib(10)
		}
	})
	b.Run("FibSlow()", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			FibSlow(10)
		}
	})
}
