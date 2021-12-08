from typing import List
from functools import lru_cache

# 1 <= rows, cols <= 50
# 1 <= k <= 10

# 你需要切披萨 k-1 次，得到 k 块披萨并送给别人。
# 请你返回确保每一块披萨包含 至少 一个苹果的切披萨方案数。由于答案可能是个很大的数字，请你返回它对 10^9 + 7 取余的结果。

# any可以前缀和优化，此处省略了
class Solution:
    def ways(self, pizza: List[str], k: int) -> int:
        r, c = len(pizza), len(pizza[0])

        @lru_cache(None)
        def dfs(x, y, k) -> int:
            # 最后一块披萨至少有一个苹果
            if not k:
                return int(any('A' in p[y:c] for p in pizza[x:r]))

            res = 0
            # 横切
            for i in range(x + 1, r):
                if any('A' in p[y:c] for p in pizza[x:i]):
                    res += dfs(i, y, k - 1)
            # 竖切
            for j in range(y + 1, c):
                if any('A' in p[y:j] for p in pizza[x:r]):
                    res += dfs(x, j, k - 1)

            return res

        return dfs(0, 0, k - 1) % int(1e9 + 7)

    def ways2(self, pizza: List[str], k: int) -> int:
        r, c = len(pizza), len(pizza[0])
        prefix = [[0] * (c + 1) for _ in range(r + 1)]
        for i in range(r):
            for j in range(c):
                prefix[i + 1][j + 1] = (
                    prefix[i][j + 1]
                    + prefix[i + 1][j]
                    - prefix[i][j]
                    + (1 if pizza[i][j] == 'A' else 0)
                )

        @lru_cache(None)
        def dfs(x, y, k):
            if not k:
                return prefix[r][c] - prefix[x][c] - prefix[r][y] + prefix[x][y] > 0
            res = 0
            for i in range(x + 1, r):
                if prefix[i][c] - prefix[x][c] - prefix[i][y] + prefix[x][y] > 0:
                    res += dfs(i, y, k - 1)
            for j in range(y + 1, c):
                if prefix[r][j] - prefix[x][j] - prefix[r][y] + prefix[x][y] > 0:
                    res += dfs(x, j, k - 1)
            return res

        return dfs(0, 0, k - 1) % (10 ** 9 + 7)


print(Solution().ways(pizza=["A..", "AAA", "..."], k=3))
# 输出：3
# 解释：上图展示了三种切披萨的方案。注意每一块披萨都至少包含一个苹果。
