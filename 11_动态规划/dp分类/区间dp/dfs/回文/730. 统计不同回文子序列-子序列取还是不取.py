from functools import lru_cache


# 给定一个字符串 s，返回 s 中不同的非空「回文子序列」个数 。
# 字符串 S 的长度将在[1, 1000]范围内。
# 每个字符 S[i] 将会是集合 {'a', 'b', 'c', 'd'} 中的某一个。

MOD = int(1e9 + 7)


class Solution:
    def countPalindromicSubsequences(self, s: str) -> int:
        """回文子序列个数 O(n^2*C^2)"""

        @lru_cache(None)
        def dfs(left: int, right: int, char: str) -> int:
            """[left,right]这一段里的不同回文子序列的个数,两端字母为char
            
            每一次取都是长度递增的 所以不会出现重复
            """
            if left >= right:
                return int(left == right and s[left] == char)
            if s[left] != char:
                return dfs(left + 1, right, char)
            if s[right] != char:
                return dfs(left, right - 1, char)
            # 子序列 取还是不取
            # 讨论取还是不取中间这段，不取就是后序dfs返回两端的 2 (奇/偶长度)，取就是 dfs(left+1,right−1,next)
            return (sum(dfs(left + 1, right - 1, char) for char in 'abcd') + 2) % MOD

        return sum(dfs(0, len(s) - 1, char) for char in 'abcd') % MOD

    def countPalindromicSubsequences2(self, s: str) -> int:
        """回文子序列个数 O(n^2*C)
        
        每一次取都是长度递增的 所以不会出现重复
        """

        @lru_cache(None)
        def dfs(left: int, right: int) -> int:
            if left >= right:
                return int(left == right)

            res = 0
            for cur in 'abcd':
                i, j = s.find(cur, left, right + 1), s.rfind(cur, left, right + 1)
                if i == -1:
                    continue
                res += 1 if i == j else 2 + dfs(i + 1, j - 1)  # 子序列 取还是不取
                res %= MOD
            return res

        return dfs(0, len(s) - 1)


print(Solution().countPalindromicSubsequences('bccb'))
print(Solution().countPalindromicSubsequences('bb'))
# 输出：6
# 解释：
# 6 个不同的非空回文子字符序列分别为：'b', 'c', 'bb', 'cc', 'bcb', 'bccb'。
# 注意：'bcb' 虽然出现两次但仅计数一次。

# 注意find范围是左闭右开，和slice一样
print('as'.find('s', 0, 1))

