from typing import List
from scipy.optimize import linear_sum_assignment


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
        boys, girls = robot, []
        for pos, limit in factory:
            girls.extend([pos] * limit)  # !每个工厂的位置重复limit次(拆点出来)
        costMatrix = [[INF] * len(girls) for _ in range(len(boys))]
        for i, pos1 in enumerate(boys):
            for j, pos2 in enumerate(girls):
                costMatrix[i][j] = abs(pos1 - pos2)  # 最大权匹配转换为最小权匹配

        res = linear_sum_assignment(costMatrix, maximize=False)
        return sum(costMatrix[i][j] for i, j in zip(*res))


print(
    Solution().minimumTotalDistance(
        robot=[0, 4] * 50,
        factory=[[2, 100] for _ in range(50)] + [[6, 100] for _ in range(50)],
    )
)
