# cost = max(chargeTimes) + k * sum(runningCosts) ，
# 其中 max(chargeTimes) 是这 k 个机器人中最大充电时间，
# sum(runningCosts) 是这 k 个机器人的运行时间之和。
# !请你返回在 不超过 budget 的前提下，你 最多 可以 连续 运行的机器人(子数组)数目为多少。
# !二分答案+st表查询区间最值

from collections import deque
from typing import Any, List


class MaxQueue(deque[Any]):
    def append(self, value: int, *metaInfo: Any) -> None:
        count = 1
        while self and self[-1][0] < value:
            count += self.pop()[-1]
        super().append([value, *metaInfo, count])

    def popleft(self) -> None:
        if not self:
            raise IndexError("popleft from empty queue")
        self[0][-1] -= 1
        if self[0][-1] == 0:
            super().popleft()

    @property
    def max(self) -> int:
        if not self:
            raise ValueError("maxQueue is empty")
        return self[0][0]


class Solution:
    def maximumRobots(self, chargeTimes: List[int], runningCosts: List[int], budget: int) -> int:
        """机器人是子数组
        注意到窗口移动时 表达式增降的单调性
        滑窗+单调队列
        """
        res, left, n = 0, 0, len(chargeTimes)
        maxQueue = MaxQueue()
        curSum = 0
        for right in range(n):
            maxQueue.append(chargeTimes[right])
            curSum += runningCosts[right]
            while left <= right and curSum * (right - left + 1) + maxQueue.max > budget:
                curSum -= runningCosts[left]
                maxQueue.popleft()
                left += 1
            res = max(res, right - left + 1)
        return res


print(Solution().maximumRobots(chargeTimes=[11, 12, 19], runningCosts=[10, 8, 7], budget=19))
