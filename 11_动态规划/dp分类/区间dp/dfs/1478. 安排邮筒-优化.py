# 一条笔直的高速公路上有 N 个村庄，每个村庄都有一个整数位置坐标，
# 不同村庄的坐标不同，现在要在其中的 P 个村庄上建立邮局。

# 请问如何安排邮局的位置可以使得每个村庄到其最近邮局的距离和最小，输出这个最小值。
from functools import lru_cache
from typing import List


class Solution:
    def minDistance(self, houses: List[int], k: int) -> int:
        """O(n^2*k)"""

        @lru_cache(None)
        def calDistance(left: int, right: int) -> int:
            """在[left,right]房子间放邮筒，返回最小距离之和。
            
            中位数，也可以记忆化处理
            """
            if left >= right:
                return 0
            return houses[right] - houses[left] + calDistance(left + 1, right - 1)

        @lru_cache(None)
        def dfs(pos: int, remain: int) -> int:
            if remain == 1:
                return calDistance(pos, len(houses) - 1)

            res = int(1e20)
            # 枚举分割点
            for i in range(pos + 1, len(houses) - remain + 2):
                res = min(res, calDistance(pos, i - 1) + dfs(i, remain - 1))
            return res

        houses.sort()
        return dfs(0, k)


# 1 <= n <= 100
# 1 <= houses[i] <= 10^4
# 1 <= k <= n
