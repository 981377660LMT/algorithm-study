from functools import lru_cache

# 1 <= S.length <= 200
# 有效排列 是对整数 {0, 1, ..., n} 的一个排列


MOD = int(1e9 + 7)

# summary
# from last to first
# 最后一个状态可以具有 [0， n] 之间的任何值
# 保持数字数量的信息尚未出现以及其中有多少比上一个数字低（或更高）就足够了。


class Solution:
    def numPermsDISequence(self, s: str) -> int:
        @lru_cache(None)
        def dfs(index: int, less: int, more: int) -> int:
            if index < 0:
                return 1

            res = 0

            if s[index] == 'I':
                for k in range(less):
                    res += dfs(index - 1, k, less + more - k - 1) % MOD
            else:
                for k in range(more):
                    res += dfs(index - 1, less + more - k - 1, k) % MOD

            return res

        n = len(s)
        return sum(dfs(n - 1, k, n - k) for k in range(n + 1)) % MOD


print(Solution().numPermsDISequence("DID"))
# 输出：5
# 解释：
# (0, 1, 2, 3) 的五个有效排列是：
# (1, 0, 3, 2)
# (2, 0, 3, 1)
# (2, 1, 3, 0)
# (3, 0, 2, 1)
# (3, 1, 2, 0)
