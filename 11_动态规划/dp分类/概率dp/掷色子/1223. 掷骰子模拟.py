from functools import lru_cache
from typing import List

# 投掷骰子时，连续 掷出数字 i 的次数不能超过 rollMax[i]
# 计算掷 n 次骰子可得到的不同点数序列的数量。
# 1 <= n <= 5000
# rollMax.length == 6
# 1 <= rollMax[i] <= 15

MOD = int(1e9 + 7)


class Solution:
    def dieSimulator(self, n: int, rollMax: List[int]) -> int:
        @lru_cache(None)
        def dfs(index: int, pre: int, count: int) -> int:
            if index == n:
                return 1

            res = 0
            for cur in range(1, 7):
                nextCount = 1 if cur != pre else count + 1
                if nextCount <= rollMax[cur - 1]:
                    res += dfs(index + 1, cur, nextCount)
                    res %= MOD
            return res

        res = 0
        for start in range(1, 7):
            res += dfs(1, start, 1)
            res %= MOD
        return res


print(Solution().dieSimulator(n=2, rollMax=[1, 1, 2, 2, 2, 3]))
# 输出：34
# 解释：我们掷 2 次骰子，如果没有约束的话，共有 6 * 6 = 36 种可能的组合。但是根据 rollMax 数组，数字 1 和 2 最多连续出现一次，所以不会出现序列 (1,1) 和 (2,2)。因此，最终答案是 36-2 = 34。

