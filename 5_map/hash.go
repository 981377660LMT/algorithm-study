// 在JDK（以JDK14为例）中Map的实现非常多，我们讲解的HashMap主要实现在 java.util 下的 HashMap 中，这是一个最简单的不考虑并发的、基于散列的Map实现。

package main

import "fmt"

func HashCode(value []byte) int32 {
	h := int32(0)
	length := int32(len(value)) >> 1
	for i := int32(0); i < length; i++ {
		h = 31*h + int32(value[i])
	}
	return h
}

func main() {
	fmt.Println(HashCode([]byte("Hello, World!")))
}
