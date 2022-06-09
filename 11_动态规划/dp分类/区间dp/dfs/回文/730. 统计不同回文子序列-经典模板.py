from functools import lru_cache

MOD = int(1e9 + 7)

# 给定一个字符串 s，返回 s 中不同的非空「回文子序列」个数 。
# 字符串 S 的长度将在[1, 1000]范围内。
# 每个字符 S[i] 将会是集合 {'a', 'b', 'c', 'd'} 中的某一个。


class Solution:
    def countPalindromicSubsequences(self, s: str) -> int:
        """回文子序列个数"""

        @lru_cache(None)
        def dfs(left: int, right: int) -> int:
            if left > right:
                return 0

            res = 0
            for newChar in 'abcd':
                i, j = s.find(newChar, left, right + 1), s.rfind(newChar, left, right + 1)
                if i == -1 or j == -1:
                    continue
                res += 1 if i == j else 2 + dfs(i + 1, j - 1)  # `边界两个回文`+`里面的子序列再带上外面这对`
            return res % MOD

        n = len(s)
        return dfs(0, n - 1)


print(Solution().countPalindromicSubsequences('bccb'))
print(Solution().countPalindromicSubsequences('bb'))
# 输出：6
# 解释：
# 6 个不同的非空回文子字符序列分别为：'b', 'c', 'bb', 'cc', 'bcb', 'bccb'。
# 注意：'bcb' 虽然出现两次但仅计数一次。

# 注意find范围是左闭右开，和slice一样
print('as'.find('s', 0, 1))

