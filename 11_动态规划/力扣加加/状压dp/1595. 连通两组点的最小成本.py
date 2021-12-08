from typing import List
from functools import lru_cache

# 如果两个组中的每个点都与另一组中的一个或多个点连接，则称这两组点是连通的
# 任意两点间的连接成本 cost 由大小为 size1 x size2 矩阵给出
# 1 <= size1, size2 <= 12
# size1 >= size2 。
# 返回连通两组点所需的最小成本。

# 总结：
# 注意这里我们状态需要遍历m(左边做出选择)，用n作为状态(右边都要被选到)
class Solution:
    def connectTwoGroups(self, cost: List[List[int]]) -> int:
        m, n = len(cost), len(cost[0])
        target = (1 << n) - 1
        minCost = [min(choose) for choose in cost]

        @lru_cache(None)
        def dfs(cur: int, state: int) -> int:
            if state == target:
                return 0

            res = 0x7FFFFFFF
            for next in range(m):
                if state & (1 << next):
                    continue
                res = min(res, minCost[cur] + dfs(cur + 1, state | (1 << next)))

            return res

        return dfs(0, 0)


print(Solution().connectTwoGroups(cost=[[15, 96], [36, 2]]))
# 输出：17
# 解释：连通两组点的最佳方法是：
# 1--A
# 2--B
# 总成本为 17 。
print(Solution().connectTwoGroups(cost=[[1, 3, 5], [4, 1, 1], [1, 5, 3]]))
# 输出：4
# 解释：连通两组点的最佳方法是：
# 1--A
# 2--B
# 2--C
# 3--A
# 最小成本为 4 。
# 请注意，虽然有多个点连接到第一组中的点 2 和第二组中的点 A ，但由于题目并不限制连接点的数目，所以只需要关心最低总成本
