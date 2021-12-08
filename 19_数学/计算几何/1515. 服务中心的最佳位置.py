# 使服务中心 到所有客户的欧几里得距离的总和最小
from typing import List

# 1 <= positions.length <= 50

# 凸函数求最值

# 梯度下降法
class Solution:
    def getMinDistSum(self, positions: List[List[int]]) -> float:
        ...


print(Solution().getMinDistSum(positions=[[0, 1], [1, 0], [1, 2], [2, 1]]))
