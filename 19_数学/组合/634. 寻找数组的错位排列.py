# n 的范围是 [1, 10^6]。

MOD = int(1e9 + 7)

# 给定一个从 1 到 n 升序排列的数组，你可以计算出总共有多少个不同的错位排列吗？
# https://leetcode-cn.com/problems/find-the-derangement-of-an-array/solution/xun-zhao-shu-zu-de-cuo-wei-pai-lie-by-lenn123/


class Solution:
    def findDerangement(self, n: int) -> int:
        dp = [0] * (n + 1)
        dp[0] = 1
        dp[1] = 0
        for i in range(2, n + 1):
            # 对前面每个数k 将n与k先交换 固定n的位置
            # k排到n 剩下有dp[i-2]种
            # k不排到n 我们将问题缩减成余下n-1个元素共有多少种错排的子问题 有dp[i-1]种
            dp[i] = (i - 1) * (dp[i - 1] + dp[i - 2])
            dp[i] %= MOD

        return dp[n]


print(Solution().findDerangement(3))
# 输出: 2
# 解释: 原始的数组为 [1,2,3]。两个错位排列的数组为 [2,3,1] 和 [3,1,2]。
