# 1 ≤ n < 2 ** 31
# 0 ≤ k < n
# 0 ≤ total < 2 ** 31
class Solution:
    def solve(self, n: int, total: int, k: int):
        """数组长为n,和为tital，相邻两个数差不大于1，求index k 的数的最大值"""
        # [X-k, X-K+1, ..... , X, X-1,..... X-N-k]
        const = n - k
        res = total + k * (k + 1) // 2 + const * (const - 1) // 2
        return res // n

