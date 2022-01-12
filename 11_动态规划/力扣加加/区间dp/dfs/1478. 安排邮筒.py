from typing import List
from functools import lru_cache

# 1 <= n <= 100
# 现需要在这条街上安排 k 个邮筒。
# 请你返回每栋房子与离它最近的邮筒之间的距离的 最小 总和。
class Solution:
    def minDistance(self, houses: List[int], k: int) -> int:
        houses.sort()

        def calDistance(left: int, right: int) -> int:
            res = 0
            while left < right:
                res += houses[right] - houses[left]
                left += 1
                right -= 1
            return res

        @lru_cache(None)
        def dfs(cur: int, remain: int) -> int:
            """
            Args:
                cur (int): [当前房屋索引]
                remain (int): [剩余邮筒]
            """
            if remain == 1:
                return calDistance(cur, len(houses) - 1)

            res = 0x7FFFFFFF
            # 枚举分割点
            for i in range(cur + 1, len(houses) - remain + 2):
                res = min(res, calDistance(cur, i - 1) + dfs(i, remain - 1))
            return res

        return dfs(0, k)


# print(Solution().minDistance(houses=[1, 4, 8, 10, 20], k=3))
# 输出：5
# 解释：将邮筒分别安放在位置 3， 9 和 20 处。
# 每个房子到最近邮筒的距离和为 |3-1| + |4-3| + |9-8| + |10-9| + |20-20| = 5 。

print(Solution().minDistance(houses=[7, 4, 6, 1], k=1))

