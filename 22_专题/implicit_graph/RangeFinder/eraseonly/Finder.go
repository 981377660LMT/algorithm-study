package erase_only_finder

type Finder interface {
	// 建立一个包含0到n-1的集合.
	Init(n int)

	// 查询元素i是否存在.
	Has(i int) bool
	// 尝试删除元素i.
	Erase(i int) bool
	// 返回大于等于i的最小元素.如果不存在,返回n.
	Prev(i int) int
	// 返回小于等于i的最大元素.如果不存在,返回-1.
	Next(i int) int
}
