# 3472. 至多 K 次操作后的最长回文子序列
# https://leetcode.cn/problems/longest-palindromic-subsequence-after-at-most-k-operations/description/
# 自序列：选或不选
# dp[i][j][k] 表示 s[i:j+1] 至多进行 k 次操作的最长回文子序列长度


from functools import lru_cache


class Solution:
    def longestPalindromicSubsequence(self, s: str, k: int) -> int:
        nums = list(map(ord, s))
        n = len(nums)

        @lru_cache(None)
        def dfs(i: int, j: int, k: int) -> int:
            if i >= j:
                return j - i + 1
            res = max(dfs(i + 1, j, k), dfs(i, j - 1, k))
            d = abs(nums[i] - nums[j])
            op = min(d, 26 - d)
            if op <= k:
                res = max(res, 2 + dfs(i + 1, j - 1, k - op))
            return res

        res = dfs(0, n - 1, k)
        dfs.cache_clear()
        return res

    def longestPalindromicSubsequence2(self, s: str, K: int) -> int:
        nums = list(map(ord, s))
        n = len(s)
        cnt = 0
        for i in range(n // 2):
            d = abs(nums[i] - nums[-1 - i])
            cnt += min(d, 26 - d)
        if cnt <= K:  # 变成回文串
            return n

        f = [
            [[0] * n for _ in range(n)] for _ in range(K + 1)
        ]  # 提高访问缓存的效率，把 k 放到第一个维度
        for k in range(K + 1):
            for i in range(n - 1, -1, -1):
                f[k][i][i] = 1
                for j in range(i + 1, n):
                    res = max(f[k][i + 1][j], f[k][i][j - 1])
                    d = abs(nums[i] - nums[j])
                    op = min(d, 26 - d)
                    if op <= k:
                        res = max(res, f[k - op][i + 1][j - 1] + 2)
                    f[k][i][j] = res
        return f[K][0][-1]
