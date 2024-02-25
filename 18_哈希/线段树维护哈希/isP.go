package main

import "fmt"

func main() {
	ptr := 0
	for i := int(1e9 + 7); ; i-- {
		if isPrime(i) {
			ptr++
			fmt.Println(i)
			if ptr == 10 {
				break
			}
		}
	}
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}
