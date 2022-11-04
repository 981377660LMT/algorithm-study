package sets

import (
	"fmt"

	"github.com/emirpasic/gods/queues/linkedlistqueue"
)

func b() {
	queue := linkedlistqueue.New()
	queue.Enqueue(1)
	queue.Enqueue(2)
	fmt.Println(queue.Peek()) // 1 true
	fmt.Println(queue.Values()...)
	fmt.Println(queue.Dequeue())
}
