package lists

// 不是线程安全的的 ArrayList

import (
	"fmt"

	"github.com/emirpasic/gods/lists/arraylist"
)

func a() {
	list := arraylist.New(0)
	fmt.Println(list.Size())

	// 增
	list.Add(1, 2, 3, 1)
	list.Insert(0, 30)
	// 删
	list.Remove(0)
	// 改
	list.Set(0, 100)
	// 查
	value, ok := list.Get(0)
	fmt.Println(value, ok)
	fmt.Println(list.Contains(1))
	fmt.Println(list.IndexOf(1))
	// 遍历
	for i := 0; i < list.Size(); i++ {
		fmt.Println(list.Get(i))
	}
	list.Each(func(index int, value interface{}) {
		fmt.Println(index, value)
	})

	// 其他
	list.Sort(func(a, b interface{}) int {
		return a.(int) - b.(int)
	})

	list.Swap(0, 1)

	fmt.Println(list.String())
	values := list.Values()
	for _, v := range values {
		fmt.Println(v)
	}

}
