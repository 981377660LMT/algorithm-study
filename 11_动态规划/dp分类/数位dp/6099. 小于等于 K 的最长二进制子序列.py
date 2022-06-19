from functools import lru_cache

# 时间复杂度O(n∗logk)

# 思路是对子序列的每个位置，决定选还是不选
# 不选直接跳过这个位置，选的话类似数位dp，要确定当前可选的上界以及判断是否存在前导0
# 最后当 s 串剩下的长度比 k 串小时，s串剩下的部分需要全选


class Solution:
    def longestSubsequence(self, s: str, k: int) -> int:
        """s 的 最长 子序列，且该子序列对应的 二进制 数字小于等于 k """

        @lru_cache(None)
        def dfs(i: int, j: int, isLimited: bool, hasLeadingZero: bool) -> int:
            """i,j表示匹配的位置 isLimited表示是否贴合上界 hasLeadingZero表示是否有前导0"""
            if n1 - i < n2 - j:
                return n1 - i
            if i == n1 or j == n2:
                return 0

            res = dfs(i + 1, j, isLimited, hasLeadingZero)  # 不选当前位置

            upper = int(target[j]) if isLimited else 1  # 选当前位置，求允许的上界
            ok = set(range(upper + 1))
            if int(s[i]) in ok:
                if hasLeadingZero and s[i] == '0':
                    res = max(res, 1 + dfs(i + 1, j, isLimited, True))
                else:
                    res = max(res, 1 + dfs(i + 1, j + 1, isLimited and (int(s[i]) == upper), False))
            return res

        target = bin(k)[2:]
        n1, n2 = len(s), len(target)
        return dfs(0, 0, True, True)


print(Solution().longestSubsequence(s="1001010", k=5))
