from bisect import bisect_right
from typing import List

# 如果你选择的工作在时间 X 结束，那么你可以立刻进行在时间 X 开始的下一份工作。
# 线性dp


class Solution:
    def jobScheduling(self, startTime: List[int], endTime: List[int], profit: List[int]) -> int:
        n = len(startTime)
        jobs = [(s, e, p) for s, e, p in zip(startTime, endTime, profit)]
        jobs.sort(key=lambda x: x[1])
        endTime = [e for _, e, _ in jobs]

        # dp[i]表示选择接乘客i为结尾时所能达到的最大盈利
        dp = [p for *_, p in jobs]
        for i in range(1, n):
            pre = bisect_right(endTime, jobs[i][0]) - 1
            if pre >= 0:
                dp[i] = max(dp[i - 1], dp[pre] + jobs[i][2])
            else:
                dp[i] = max(dp[i - 1], jobs[i][2])

        return dp[-1]


print(Solution().jobScheduling([1, 2, 3, 3], [3, 4, 5, 6], [50, 10, 40, 70]))
