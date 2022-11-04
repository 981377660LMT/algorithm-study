// 循环链表

package lists

import (
	"fmt"

	linkedList "github.com/emirpasic/gods/lists/doublylinkedlist"
)

func b() {
	linkedList := linkedList.New()
	linkedList.Add(1, 2, 3, 4, 5)
	fmt.Println(linkedList.String())

	// 入队
	linkedList.Append(6)
	linkedList.Prepend(0)
	fmt.Println(linkedList.String())

	// 出队
	linkedList.Remove(0)
	linkedList.Remove(linkedList.Size() - 1)
	fmt.Println(linkedList.String())

	// 遍历
	linkedList.Each(func(index int, value interface{}) {
		if value == 3 {
			fmt.Println(index)
		}
	})

	json, err := linkedList.ToJSON()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(json))
}
