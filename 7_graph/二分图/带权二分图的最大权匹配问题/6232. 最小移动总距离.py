from typing import List
from KM算法模板 import KM


INF = int(1e18)


class Solution:
    def minimumTotalDistance(self, robot: List[int], factory: List[List[int]]) -> int:
        """带权图最大权匹配

        Args:
            robot (List[int]): 机器人的位置
            factory (List[List[int]]): 每个工厂的(位置,可以修理的机器人个数)

        Returns:
            int: _description_ 请你返回所有机器人移动的最小总距离。测试数据保证所有机器人都可以被维修。
        """
        girls, boys = [], robot
        for pos, limit in factory:
            girls.extend([pos] * limit)  # !每个工厂的位置重复limit次(拆点出来)
        costMatrix = [[0] * len(girls) for _ in range(len(boys))]
        for i, pos1 in enumerate(boys):
            for j, pos2 in enumerate(girls):
                costMatrix[i][j] = -abs(pos1 - pos2)  # 最大权匹配转换为最小权匹配

        return -KM(costMatrix)[0]


print(Solution().minimumTotalDistance(robot=[0, 4, 6], factory=[[2, 2], [6, 2]]))
