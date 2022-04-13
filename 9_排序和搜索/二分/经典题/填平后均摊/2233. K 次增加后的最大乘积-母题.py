# 假设此题k很大 不可用堆模拟
from functools import reduce
from itertools import accumulate
from typing import List


MOD = int(1e9 + 7)
INF = int(1e20)

# 6039. K 次增加后的最大乘积
# k次加1操作，让最小值最大化
# 填平+均摊

# k<=1e9

MOD = int(1e9 + 7)


class Solution:
    def maximumProduct(self, nums: List[int], k: int) -> int:
        def maximizeMinValue(nums: List[int], delta: int) -> List[int]:
            """delta次加1操作，让最小值最大化，返回操作后的数组"""
            n = len(nums)
            copy = nums[:]
            nums = sorted(nums)
            preSum = [0] + list(accumulate(nums))
            nums = [0] + nums

            # 最右二分求最后能和哪个数齐平
            left, right = 0, n
            while left <= right:
                mid = (left + right) >> 1
                diff = mid * nums[mid] - preSum[mid]
                if diff <= delta:
                    left = mid + 1
                else:
                    right = mid - 1

            min_ = nums[right]
            overflow = delta - (right * nums[right] - preSum[right])
            div, mod = 0, 0
            if right:
                div, mod = divmod(overflow, right)
            min_ += div

            for i in range(n):
                if copy[i] < min_ + int(mod > 0):
                    copy[i] = min_ + int(mod > 0)
                    mod -= 1

            return copy

        nums = maximizeMinValue(nums, k)
        return reduce(lambda pre, cur: pre * cur % MOD, nums, 1)
        res = 1
        for num in nums:
            res *= num
            res %= MOD
        return res


print(Solution().maximumProduct(nums=[6, 3, 3, 2], k=2))
