from typing import List
from MinCostMaxFlow import MinCostMaxFlowEK

INF = int(1e18)


class Solution:
    def minimumTotalDistance(self, robot: List[int], factory: List[List[int]]) -> int:
        """最小费用最大流

        Args:
            robot (List[int]): 机器人的位置
            factory (List[List[int]]): 每个工厂的(位置,可以修理的机器人个数)

        Returns:
            int: _description_ 请你返回所有机器人移动的最小总距离。测试数据保证所有机器人都可以被维修。
        """
        n, m = len(robot), len(factory)
        STRAT, END = n + m, n + m + 1
        mcmf = MinCostMaxFlowEK(n + m + 2, STRAT, END)
        for i in range(n):
            mcmf.addEdge(STRAT, i, 1, 0)
        for i in range(n):
            for j in range(m):
                mcmf.addEdge(i, n + j, 1, abs(robot[i] - factory[j][0]))
        for i in range(m):
            mcmf.addEdge(n + i, END, factory[i][1], 0)
        return mcmf.work()[1]


print(Solution().minimumTotalDistance(robot=[0, 4, 6], factory=[[2, 2], [6, 2]]))
