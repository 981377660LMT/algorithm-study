from typing import List
from math import ceil


class Solution:
    def eliminateMaximum(self, dist: List[int], speed: List[int]) -> int:
        n = len(dist)
        time = sorted(ceil(l / v) for l, v in zip(dist, speed))

        for i, t in enumerate(time):
            if t <= i:
                return i

        return n


print(Solution().eliminateMaximum(dist=[1, 1, 2, 3], speed=[1, 1, 1, 1]))
# 输出：1
# 解释：
# 第 0 分钟开始时，怪物的距离是 [1,1,2,3]，你消灭了第一个怪物。
# 第 1 分钟开始时，怪物的距离是 [X,0,1,2]，你输掉了游戏。
# 你只能消灭 1 个怪物。
