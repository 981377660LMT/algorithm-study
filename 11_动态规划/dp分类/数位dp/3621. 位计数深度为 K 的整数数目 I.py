# 3621. 位计数深度为 K 的整数数目 I
# https://leetcode.cn/problems/number-of-integers-with-popcount-depth-equal-to-k-i/description/
#
# 给你两个整数 n 和 k。
#
# 对于任意正整数 x，定义以下序列：
# p0 = x
# pi+1 = popcount(pi)，对于所有 i >= 0，其中 popcount(y) 是 y 的二进制表示中 1 的数量。
# 这个序列最终会达到值 1。
#
# x 的 popcount-depth （位计数深度）定义为使得 pd = 1 的 最小 整数 d >= 0。
#
# 例如，如果 x = 7（二进制表示 "111"）。那么，序列是：7 → 3 → 2 → 1，所以 7 的 popcount-depth 是 3。
#
# 你的任务是确定范围 [1, n] 中 popcount-depth 恰好 等于 k 的整数数量。
#
# 返回这些整数的数量。


from functools import lru_cache


class Solution:
    def popcountDepth(self, n: int, k: int) -> int:
        @lru_cache(None)
        def dfs(index: int, isLimit: int, remain1: int) -> int:
            if index == m:
                return 1 if remain1 == 0 else 0
            upper = nums[index] if isLimit else 1
            res = 0
            for d in range(min(upper, remain1) + 1):
                res += dfs(index + 1, isLimit and (d == upper), remain1 - d)
            return res

        if k == 0:
            return 1

        nums = list(map(int, bin(n)[2:]))
        m = len(nums)
        if k == 1:
            return m - 1

        depth = [0] * (m + 1)
        res = 0
        for i in range(1, m + 1):
            depth[i] = depth[i.bit_count()] + 1
            if depth[i] == k:
                res += dfs(0, True, i)
        dfs.cache_clear()
        return res
