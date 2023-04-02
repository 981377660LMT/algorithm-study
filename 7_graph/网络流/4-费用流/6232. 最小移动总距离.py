from typing import List
from MinCostMaxFlow import MinCostMaxFlowDinic

INF = int(1e18)


class Solution:
    def minimumTotalDistance(self, robot: List[int], factory: List[List[int]]) -> int:
        """最小费用最大流 两两暴力连边

        Args:
            robot (List[int]): 机器人的位置
            factory (List[List[int]]): 每个工厂的(位置,可以修理的机器人个数)

        Returns:
            int: _description_ 请你返回所有机器人移动的最小总距离。测试数据保证所有机器人都可以被维修。
        """
        n, m = len(robot), len(factory)
        STRAT, END = n + m, n + m + 1
        mcmf = MinCostMaxFlowDinic(n + m + 2, STRAT, END)
        for i in range(n):
            mcmf.addEdge(STRAT, i, 1, 0)
        for i in range(n):
            for j in range(m):
                mcmf.addEdge(i, n + j, 1, abs(robot[i] - factory[j][0]))
        for i in range(m):
            mcmf.addEdge(n + i, END, factory[i][1], 0)
        return mcmf.work()[1]

    def minimumTotalDistance2(self, robot: List[int], factory: List[List[int]]) -> int:
        """利用排序优化建图 边数从O(nm)降到O(n+m)

        !将工厂和机器人按照坐标排序,只连接X轴上相邻的两点,
        费用为两点间距离,容量为n(或者大于n的任意数值),总边数O(n+m)
        """
        n, m = len(robot), len(factory)
        STRAT, END = n + m, n + m + 1
        positions = []  # (pos,vertex)
        mcmf = MinCostMaxFlowDinic(n + m + 2, STRAT, END)
        for i, r in enumerate(robot):  # 源点到机器人
            positions.append((r, i))
            mcmf.addEdge(STRAT, i, 1, 0)
        for i, (f, c) in enumerate(factory):  # 工厂到汇点
            positions.append((f, n + i))
            mcmf.addEdge(n + i, END, c, 0)

        positions.sort(key=lambda x: x[0])  # 按照坐标排序
        for i in range(len(positions) - 1):  # !连接X轴上相邻的两点 容量无穷大(等价于相互建边)
            pos1, v1 = positions[i]
            pos2, v2 = positions[i + 1]
            mcmf.addEdge(v1, v2, INF, pos2 - pos1)
            mcmf.addEdge(v2, v1, INF, pos2 - pos1)
        return mcmf.work()[1]


print(Solution().minimumTotalDistance(robot=[0, 4, 6], factory=[[2, 2], [6, 2]]))
print(Solution().minimumTotalDistance2(robot=[0, 4, 6], factory=[[2, 2], [6, 2]]))
