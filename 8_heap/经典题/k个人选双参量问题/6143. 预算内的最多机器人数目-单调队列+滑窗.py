# cost = max(chargeTimes) + k * sum(runningCosts) ，
# 其中 max(chargeTimes) 是这 k 个机器人中最大充电时间，
# sum(runningCosts) 是这 k 个机器人的运行时间之和。
# !请你返回在 不超过 budget 的前提下，你 最多 可以 连续 运行的机器人(子数组)数目为多少。
# !排序+堆 `进来的越没用，出去的最没用`


from heapq import heappop, heappush
from typing import List


class Solution:
    def maximumRobots2(self, chargeTimes: List[int], runningCosts: List[int], budget: int) -> int:
        """
        机器人是子序列

        排序+堆维护 双变量制约的题
        `乘参量 1 用排序`,`参量 2 用堆维护`,出堆入堆形成抗衡,同时更新 res
        """
        robots = sorted(zip(chargeTimes, runningCosts), key=lambda x: x[0])
        pq = []  # 维护最大值
        res, curSum = 0, 0
        for max_, num in robots:
            curSum += num
            heappush(pq, -num)
            while pq and max_ + curSum * len(pq) > budget:
                curSum += heappop(pq)
            res = max(res, len(pq))
        return res


print(Solution().maximumRobots2(chargeTimes=[11, 12, 19], runningCosts=[10, 8, 7], budget=19))
