from typing import List
from functools import lru_cache

# 如果两个组中的每个点都与另一组中的一个或多个点连接，则称这两组点是连通的
# 任意两点间的连接成本 cost 由大小为 size1 x size2 矩阵给出
# 1 <= size1, size2 <= 12
# size1 >= size2 。
# 返回连通两组点所需的最小成本。

# 总结：
# 1.因为已知第二组点的数量较少，所以对第二组点的连通状态进行状态压缩，然后依次处理第一组中的点即可。
# 2.左边选完后，对右边没选的，用最小值相连


class Solution:
    def connectTwoGroups(self, cost: List[List[int]]) -> int:
        n, m = len(cost), len(cost[0])
        # 右边每个点连接的最小费用
        minCost = [min(c) for c in zip(*cost)]

        @lru_cache(None)
        def dfs(index: int, state: int) -> int:
            if index == n:
                remainCost = 0
                for j in range(m):
                    if not state & (1 << j):
                        remainCost += minCost[j]
                return remainCost

            res = int(1e9)
            for next in range(m):
                res = min(res, cost[index][next] + dfs(index + 1, state | (1 << next)))

            return res

        return dfs(0, 0)


print(Solution().connectTwoGroups(cost=[[15, 96], [36, 2]]))
# 输出：17
# 解释：连通两组点的最佳方法是：
# 1--A
# 2--B
# 总成本为 17 。
