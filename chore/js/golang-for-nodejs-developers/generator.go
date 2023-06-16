// 1. channel 实现 generator, 类似 for of 的效果
// 2. 迭代器实现 generator. 性能更好.

package main

import "fmt"

func Generator() chan string {
	c := make(chan string)
	go func() {
		c <- "hello"
		c <- "world"
		close(c)
	}()
	return c
}

func GeneratorFunc() func() (string, bool) {
	s := []string{"hello", "world"}
	i := -1
	return func() (string, bool) {
		i++
		if i >= len(s) {
			return "", false
		}
		return s[i], true
	}
}

func main() {
	gen := Generator()
	for {
		value, more := <-gen
		fmt.Println(value, more)
		if !more {
			break
		}
	}

	// alternatively
	for value := range Generator() {
		fmt.Println(value)
	}

	////////////////////////////////////////////////

	// alternatively
	genfn := GeneratorFunc()
	for {
		value, more := genfn()
		fmt.Println(value, more)

		if !more {
			break
		}
	}
}
