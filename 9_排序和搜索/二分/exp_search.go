package main

import "fmt"

func main() {
	// --- 用例 1: 在有序数组中查找小于等于目标值的最后一个元素 ---
	nums := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19}
	target := 12
	fmt.Printf("数组: %v, 目标: <= %d\n", nums, target)

	// 我们要找最大的索引 i，使得 nums[i] <= 12
	// 起始点设为 -1 (假设数组可能全都不满足) 或者 0 (如果确定 nums[0] <= target)
	// 这里为了安全，我们假设 nums[0] 满足条件，从 0 开始搜
	if len(nums) > 0 && nums[0] <= target {
		index := ExpSearch(func(i int) bool {
			// 注意：必须处理数组越界，因为 ExpSearch 可能会探测超出数组长度的索引
			if i >= len(nums) {
				return false
			}
			return nums[i] <= target
		}, 0)
		fmt.Printf("找到索引: %d, 值: %d\n", index, nums[index])
	} else {
		fmt.Println("没有找到符合条件的元素")
	}

	fmt.Println("-------------------")

	// --- 用例 2: 计算整数平方根 (无界搜索) ---
	// 寻找最大的整数 x，使得 x * x <= N
	N := 2000
	fmt.Printf("计算 %d 的整数平方根\n", N)

	// 从 0 开始搜索，因为 0*0 <= 2000 肯定成立
	sqrt := ExpSearch(func(x int) bool {
		return x*x <= N
	}, 0)

	fmt.Printf("结果: %d (验证: %d*%d=%d, %d*%d=%d)\n",
		sqrt, sqrt, sqrt, sqrt*sqrt, sqrt+1, sqrt+1, (sqrt+1)*(sqrt+1))

	fmt.Println("-------------------")

	// --- 用例 3: 模拟未知长度的数据流/API ---
	// 假设有一个函数 API，我们不知道它什么时候返回 false，只知道它在某个点之后会一直失败
	// 比如：寻找某个服务支持的最大并发连接数
	maxCapacity := 532 // 假设这是真实限制，但我们不知道
	fmt.Printf("寻找最大容量 (真实值: %d)\n", maxCapacity)

	checkConnection := func(n int) bool {
		// 模拟 API 调用
		return n <= maxCapacity
	}

	// 从 1 开始尝试
	capacity := ExpSearch(checkConnection, 1)
	fmt.Printf("探测到的最大容量: %d\n", capacity)
}

// ExpSearch 执行指数搜索（倍增搜索）。
// 寻找满足 check(x) 为 true 的最大索引 x。
// !check 函数必须满足单调性，即 [true, true, ..., true, false, false, ...]
//
// 参数:
//
//	check: 判定函数，输入索引，返回是否满足条件。
//	ok:    起始索引，必须保证 check(ok) 为 true。
//
// 返回值:
//
//	满足 check 条件的最大的 int 值。
//
// 时间复杂度: O(log i)，其中 i 是返回值与 ok 的差值。
func ExpSearch(check func(int) bool, ok int) int {
	d := 1
	for check(ok + d) {
		ok += d
		d += d
	}
	ng := ok + d
	for ok+1 < ng {
		mid := (ok + ng) / 2
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}
