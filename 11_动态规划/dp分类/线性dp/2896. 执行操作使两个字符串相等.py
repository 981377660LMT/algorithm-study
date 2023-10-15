# 2896. 执行操作使两个字符串相等
# https://leetcode.cn/problems/apply-operations-to-make-two-strings-equal/description/
# 给你两个下标从 0 开始的二进制字符串 s1 和 s2 ，两个字符串的长度都是 n ，再给你一个正整数 x 。
# 你可以对字符串 s1 执行以下操作 任意次 ：
# !选择两个下标 i 和 j ，将 s1[i] 和 s1[j] 都反转，操作的代价为 x 。
# !选择满足 i < n - 1 的下标 i ，反转 s1[i] 和 s1[i + 1] ，操作的代价为 1 。
# !请你返回使字符串 s1 和 s2 相等的 最小 操作代价之和，如果无法让二者相等，返回 -1 。
# 注意 ，反转字符的意思是将 0 变成 1 ，或者 1 变成 0 。
#
#
# !dp，要意识到无后效性


from functools import lru_cache

INF = int(1e18)


def min(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def minOperations(self, s1: str, s2: str, x: int) -> int:
        """操作1相当于花费x/2翻转一个位置,操作2相当于花费pos[i]-pos[i-1]翻转两个位置.
        则 `dp[i] = min(dp[i-1]+x,dp[i-2]+2*(pos[i]-pos[i-1]))`
        """
        pos = [i for i, (a, b) in enumerate(zip(s1, s2)) if a != b]  # 待翻转的位置
        k = len(pos)
        if k & 1:
            return -1

        dp = [INF] * (k + 1)
        dp[0] = 0
        for i in range(1, k + 1):
            dp[i] = min(dp[i - 1] + x, dp[i - 2] + 2 * (pos[i - 1] - pos[i - 2]))
        return dp[k] // 2

    def minOperations2(self, s1: str, s2: str, x: int) -> int:
        """O(n^2)记忆化dfs."""

        @lru_cache(None)
        def dfs(index: int, flip1Count: int, prevFlip2: bool) -> int:
            if index == n:
                return 0 if flip1Count == 0 else INF

            cur = nums1[index] ^ prevFlip2
            target = nums2[index]
            if cur == target:
                return dfs(index + 1, flip1Count, False)

            # 反转1
            res = dfs(index + 1, flip1Count + 1, False) + x
            if flip1Count > 0:
                res = min(res, dfs(index + 1, flip1Count - 1, False))
            # 反转2
            if index < n - 1:
                res = min(res, dfs(index + 1, flip1Count, True) + 1)
            return res

        nums1 = [int(i) for i in s1]
        nums2 = [int(i) for i in s2]
        n = len(s1)
        res = dfs(0, 0, False)
        dfs.cache_clear()
        return res if res < INF else -1
