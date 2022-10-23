from bisect import bisect_right
from typing import List

# 时间上出现重叠的 2 份工作不能同时进行。
# 如果你选择的工作在时间 X 结束，那么你可以立刻进行在时间 X 开始的下一份工作。
# 请你计算并返回可以获得的最大报酬

# 线性dp 选还是不选
# dp[i]=max(dp[i-1],score[i]+dp[pre])


class Solution:
    def jobScheduling(self, startTime: List[int], endTime: List[int], profit: List[int]) -> int:
        n = len(startTime)
        jobs = [(s, e, p) for s, e, p in zip(startTime, endTime, profit)]
        jobs.sort(key=lambda x: x[1])

        dp = [0] * (n + 1)  # 前i个工作的最大收益
        for i in range(n):
            dp[i + 1] = dp[i]  # 不选
            start, _, score = jobs[i]  # 选
            prePos = bisect_right(jobs, start, key=lambda x: x[1]) - 1
            dp[i + 1] = max(dp[i + 1], score + dp[prePos + 1])
        return dp[-1]


print(Solution().jobScheduling([1, 2, 3, 3], [3, 4, 5, 6], [50, 10, 40, 70]))
print(
    Solution().jobScheduling(
        startTime=[1, 2, 3, 4, 6], endTime=[3, 5, 10, 6, 9], profit=[20, 20, 100, 70, 60]
    )
)
