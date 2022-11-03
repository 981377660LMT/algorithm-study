package main

import (
	"fmt"

	"github.com/emirpasic/gods/queues/circularbuffer"
)

func main() {
	queue := circularbuffer.New(100)
	queue.Enqueue(1)
	queue.Enqueue(2)
	fmt.Println(queue.Peek()) // 1 true
	fmt.Println(queue.Values())
}
