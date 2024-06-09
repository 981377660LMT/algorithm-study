// # 相邻不等元素的对数<=k，求最长子序列.
// class Solution:
//     def maximumLength(self, nums: List[int], k: int) -> int:
//         @lru_cache(None)
//         def dfs(index: int, pre: int, count: int) -> int:
//             if count > k:
//                 return -INF
//             if index == n:
//                 return 0

//             res = 0
//             cur = nums[index]
//             bad = pre != -1 and cur != pre
//             # 选
//             res = max2(res, dfs(index + 1, cur, count + bad) + 1)
//             # 不选
//             res = max2(res, dfs(index + 1, pre, count))
//             return res

//         n = len(nums)
//         res = dfs(0, -1, 0)
//         dfs.cache_clear()
//         return res

impl Solution {
    pub fn maximum_length(nums: Vec<i32>, k: i32) -> i32 {}
}
