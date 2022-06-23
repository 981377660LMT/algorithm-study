from bisect import bisect_right
from typing import List

# 如果你选择的工作在时间 X 结束，那么你可以立刻进行在时间 X 开始的下一份工作。
# 时间上出现重叠的 2 份工作不能同时进行。
# 线性dp


class Solution:
    def jobScheduling(self, startTime: List[int], endTime: List[int], profit: List[int]) -> int:
        n = len(startTime)
        jobs = [(s, e, p) for s, e, p in zip(startTime, endTime, profit)]
        jobs.sort(key=lambda x: x[1])

        # dp[i]表示选择接乘客i为结尾时所能达到的最大盈利
        dp = [score for *_, score in jobs]
        for i in range(1, n):
            start, _, score = jobs[i]
            pre = bisect_right(jobs, start, key=lambda x: x[1]) - 1
            if pre >= 0:
                dp[i] = max(dp[i - 1], dp[pre] + score)
            else:
                dp[i] = max(dp[i - 1], score)

        return dp[-1]


print(Solution().jobScheduling([1, 2, 3, 3], [3, 4, 5, 6], [50, 10, 40, 70]))
