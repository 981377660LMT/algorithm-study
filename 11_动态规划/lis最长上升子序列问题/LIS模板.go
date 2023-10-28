// """
// 贪心 + 二分查找
// LIS[i]表示长度为 i+1 的子序列尾部元素的值
// 每次遍历到一个新元素,用二分查找法找到第一个大于等于它的元素,然后更新LIS
// """
// # LIS模板

// from typing import List, Tuple
// from bisect import bisect_left, bisect_right

// def LIS(nums: List[int], isStrict=True) -> int:
//     """求LIS长度"""
//     n = len(nums)
//     if n <= 1:
//         return n

//     lis = []  # lis[i] 表示长度为 i 的上升子序列的最小末尾值
//     f = bisect_left if isStrict else bisect_right
//     for i in range(n):
//         pos = f(lis, nums[i])
//         if pos == len(lis):
//             lis.append(nums[i])
//         else:
//             lis[pos] = nums[i]

//     return len(lis)

// def LISDp(nums: List[int], isStrict=True) -> List[int]:
//     """求以每个位置为结尾的LIS长度(包括自身)"""
//     if not nums:
//         return []
//     n = len(nums)
//     res = [0] * n
//     lis = []
//     f = bisect_left if isStrict else bisect_right
//     for i in range(n):
//         pos = f(lis, nums[i])
//         if pos == len(lis):
//             lis.append(nums[i])
//             res[i] = len(lis)
//         else:
//             lis[pos] = nums[i]
//             res[i] = pos + 1
//     return res

// def getLIS(nums: List[int], isStrict=True) -> Tuple[List[int], List[int]]:
//     """求LIS 返回(LIS,LIS的组成下标)"""
//     n = len(nums)

//     lis = []  # lis[i] 表示长度为 i 的上升子序列的最小末尾值
//     dpIndex = [0] * n  # 每个元素对应的LIS长度
//     f = bisect_left if isStrict else bisect_right
//     for i in range(n):
//         pos = f(lis, nums[i])
//         if pos == len(lis):
//             lis.append(nums[i])
//         else:
//             lis[pos] = nums[i]
//         dpIndex[i] = pos

//     res, resIndex = [], []
//     j = len(lis) - 1
//     for i in range(n - 1, -1, -1):
//         if dpIndex[i] == j:
//             res.append(nums[i])
//             resIndex.append(i)
//             j -= 1
//     return res[::-1], resIndex[::-1]

// def LISMaxSum(nums: List[int], isStrict=True) -> List[int]:
//     """求以每个位置为结尾的LIS最大和(包括自身)"""

//     def max(a: int, b: int) -> int:
//         return a if a > b else b

//     class BITPrefixMax:
//         __slots__ = ("_max", "_tree")

//         def __init__(self, max: int):
//             self._max = max
//             self._tree = dict()

//         def set(self, index: int, value: int) -> None:
//             index += 1
//             while index <= self._max:
//                 self._tree[index] = max(self._tree.get(index, 0), value)
//                 index += index & -index

//         def query(self, end: int) -> int:
//             """Query max of [0, end)."""
//             if end > self._max:
//                 end = self._max
//             res = 0
//             while end > 0:
//                 res = max(res, self._tree.get(end, 0))
//                 end -= end & -end
//             return res

//     n = len(nums)
//     if n <= 1:
//         return nums[:]
//     max_ = 0
//     for v in nums:
//         max_ = max(max_, v)
//     dp = BITPrefixMax(max_ + 5)
//     res = [0] * n
//     for i, v in enumerate(nums):
//         preMax = dp.query(v) if isStrict else dp.query(v + 1)
//         cur = preMax + v
//         res[i] = cur
//         dp.set(v, cur)
//     return res

// # // LIS 方案数 O(nlogn)
// # // 原理见下面这题官方题解的方法二
// # // LC673 https://leetcode-cn.com/problems/number-of-longest-increasing-subsequence/
// # cntLis := func(a []int) int {
// # 	g := [][]int{}   // 保留所有历史信息
// # 	cnt := [][]int{} // 个数前缀和
// # 	for _, v := range a {
// # 		p := sort.Search(len(g), func(i int) bool { return g[i][len(g[i])-1] >= v })
// # 		c := 1
// # 		if p > 0 {
// # 			// 根据 g[p-1] 来计算 cnt
// # 			i := sort.Search(len(g[p-1]), func(i int) bool { return g[p-1][i] < v })
// # 			c = cnt[p-1][len(cnt[p-1])-1] - cnt[p-1][i]
// # 		}
// # 		if p == len(g) {
// # 			g = append(g, []int{v})
// # 			cnt = append(cnt, []int{0, c})
// # 		} else {
// # 			g[p] = append(g[p], v)
// # 			cnt[p] = append(cnt[p], cnt[p][len(cnt[p])-1]+c)
// # 		}
// # 	}
// # 	c := cnt[len(cnt)-1]
// # 	return c[len(c)-1]
// # }
