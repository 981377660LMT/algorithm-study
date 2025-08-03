# https://leetcode.cn/problems/earliest-finish-time-for-land-and-water-rides-ii/solutions/3740884/on-jian-ji-xie-fa-pythonjavacgo-by-endle-xgjk/
#
# 给定两组游乐设施：陆地游乐设施和水上游乐设施。每组设施有两个属性：
# 1. `startTime[i]`：第 i 个设施的最早开放时间；
# 2. `duration[i]`：第 i 个设施的持续时间。
#
# 游客必须从每组中选择一个设施，完成两个设施的顺序不限。
# 设施可以在开放时间开始，也可以在之后任意时间开始。完成一个设施后，游客可以立即开始另一个设施（如果它已经开放），或者等待它开放。
#
# 目标是求完成两个设施的最早可能时间。
#
# 对于两组设施，分别选择一个设施，使得完成两个设施的总时间最小。


from typing import List


class Solution:
    def earliestFinishTime(
        self,
        landStartTime: List[int],
        landDuration: List[int],
        waterStartTime: List[int],
        waterDuration: List[int],
    ) -> int:
        def solve(
            startTime1: List[int], duration1: List[int], startTime2: List[int], duration2: List[int]
        ) -> int:
            minFinish1 = min(s + d for s, d in zip(startTime1, duration1))
            return min(max(minFinish1, s) + d for s, d in zip(startTime2, duration2))

        res1 = solve(landStartTime, landDuration, waterStartTime, waterDuration)
        res2 = solve(waterStartTime, waterDuration, landStartTime, landDuration)
        return min(res1, res2)
