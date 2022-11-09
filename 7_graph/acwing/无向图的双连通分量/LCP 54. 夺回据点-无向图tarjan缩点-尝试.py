from typing import List
from collections import defaultdict
from Tarjan import Tarjan


# 为了防止魔物暴动，勇者在每一次夺回据点后（包括花费资源夺回据点后），
# 需要保证剩余的所有魔物据点之间是相连通的（不经过「已夺回据点」）。


# Tarjan缩点
# 先用 Tarjan 算法找出割点，去掉这些点会剩下若干个连通块。
# 抛弃掉同时与多个割点相连的连通块。
# 求出剩余的连通块中的最小权值。
# 如果仅有一个连通块，答案就是这个最小权值；否则，答案为所有最小权值之和减去它们的最大值。
class Solution:
    def minimumCost(self, cost: List[int], roads: List[List[int]]) -> int:
        n = len(cost)
        adjMap = defaultdict(set)
        for u, v in roads:
            adjMap[u].add(v)
            adjMap[v].add(u)

        # 找VBCC和割点
        VBCCId, VBCCGroup, VBCCIdByNode = Tarjan.getVBCC(n, adjMap)
        cuttingPoints = set(i for i in range(n) if len(VBCCIdByNode[i]) > 1)

        # 统计连通分量里包含几个原图的割点，不能选连了两个以上个点的分量
        counter = [sum(node in cuttingPoints for node in VBCCGroup[i]) for i in range(VBCCId)]
        goodGroups = [i for i in range(VBCCId) if counter[i] <= 1]

        # 不能选割点
        costs = [min(cost[v] for v in VBCCGroup[k] if v not in cuttingPoints) for k in goodGroups]
        return costs[0] if len(costs) == 1 else sum(costs) - max(costs)


print(
    Solution().minimumCost(
        cost=[1, 2, 3, 4, 5, 6], roads=[[0, 1], [0, 2], [1, 3], [2, 3], [1, 2], [2, 4], [2, 5]]
    )
)

print(Solution().minimumCost(cost=[3, 2, 1, 4], roads=[[0, 2], [2, 3], [3, 1]]))


# print(Solution().minimumCost(cost=[0, 1, 2, 3], roads=[[0, 1], [1, 2], [2, 3], [0, 3]]))
print(
    Solution().minimumCost(
        cost=[9, 2, 3, 4, 5, 6, 7],
        roads=[[1, 2], [1, 3], [2, 3], [3, 6], [6, 0], [0, 3], [4, 2], [2, 5], [4, 5]],
    )
)
