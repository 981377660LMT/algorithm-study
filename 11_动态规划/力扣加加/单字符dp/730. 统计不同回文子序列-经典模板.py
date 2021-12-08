from functools import lru_cache

MOD = int(1e9 + 7)

# 字符串 S 的长度将在[1, 1000]范围内。
# 每个字符 S[i] 将会是集合 {'a', 'b', 'c', 'd'} 中的某一个。
class Solution:
    def countPalindromicSubsequences(self, s: str) -> int:
        @lru_cache(maxsize=None)
        def dfs(left: int, right: int) -> int:
            if left >= right:
                return 0

            res = 0
            for char in 'abcd':
                i, j = s.find(char, left, right), s.rfind(char, left, right)
                if i == -1 or j == -1:
                    continue
                res += 1 if i == j else 2 + dfs(i + 1, j)
            return res % MOD

        return dfs(0, len(s))


print(Solution().countPalindromicSubsequences('bccb'))
# 输出：6
# 解释：
# 6 个不同的非空回文子字符序列分别为：'b', 'c', 'bb', 'cc', 'bcb', 'bccb'。
# 注意：'bcb' 虽然出现两次但仅计数一次。

# 注意find范围是左闭右开，和slice一样
print('as'.find('s', 0, 1))

