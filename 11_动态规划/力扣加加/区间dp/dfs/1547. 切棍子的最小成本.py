from typing import List
from functools import lru_cache


# 有一根长度为 n 个单位的木棍，棍上从 0 到 n 标记了若干位置
# 其中 cuts[i] 表示你需要将棍子切开的位置。
# 你可以按顺序完成切割，也可以根据需要更改切割的顺序。
# 每次切割的成本都是当前要切割的棍子的长度

# 2 <= n <= 10^6

# 1000. 合并石头的最低成本.py
class Solution:
    def minCost(self, n: int, cuts: List[int]) -> int:
        cuts = [0] + cuts + [n]

        @lru_cache(None)
        def dfs(left: int, right: int) -> int:
            if left + 1 >= right:
                return 0
            res = 0x7FFFFFFF
            for i in range(left + 1, right):
                # 这一段的长度cuts[right] - cuts[left]
                res = min(res, cuts[right] - cuts[left] + dfs(left, i) + dfs(i, right))
            return res

        return dfs(0, len(cuts) - 1)


print(Solution().minCost(n=7, cuts=[1, 3, 4, 5]))
