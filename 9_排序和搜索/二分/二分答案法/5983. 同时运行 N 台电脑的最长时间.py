from typing import List

# 1 <= n <= batteries.length <= 105
# 1 <= batteries[i] <= 10^9
# 优先队列模拟不太好做，因为 1 <= batteries[i] <= 10^9


# !换一种问法:五个数，每次选择四个数减4，最多能操作几次
# 相当于电脑数为4，电池数为5的特殊情况
class Solution:
    def maxRunTime(self, n: int, batteries: List[int]) -> int:
        batteries = sorted(batteries, reverse=True)
        # 备用电池
        spareSum = sum(batteries[n:])

        def check(needTime: int) -> bool:
            """ "备用电池能否提供足够的储备"""
            needSpare = 0
            for i in range(n):
                supply = batteries[i]
                if supply < needTime:
                    needSpare += needTime - supply
            return needSpare <= spareSum

        left, right = 0, sum(batteries)
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1
        return right  # 最右二分


print(Solution().maxRunTime(n=2, batteries=[3, 3, 3]))
print(Solution().maxRunTime(n=2, batteries=[1, 1, 1, 1]))
print(Solution().maxRunTime(n=2, batteries=[1, 1, 1, 1]))
