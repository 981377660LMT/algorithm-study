# 假设此题k很大 不可用堆模拟
from itertools import accumulate
from typing import List
from heapq import heapify, heapreplace

MOD = int(1e9 + 7)
INF = int(1e20)

# 6039. K 次增加后的最大乘积
# k次加1操作，让最小值最大化
# 填平+均摊

# k<=1e9

MOD = int(1e9 + 7)


class Solution:
    def maximumProduct(self, nums: List[int], k: int) -> int:
        """k次加1操作，让最小值最大化"""
        n = len(nums)
        arr = sorted(nums)
        preSum = [0] + list(accumulate(arr))
        arr = [0] + arr

        # 最右二分求最后能和哪个数齐平
        left, right = 0, n
        while left <= right:
            mid = (left + right) >> 1
            diff = mid * arr[mid] - preSum[mid]
            if diff <= k:
                left = mid + 1
            else:
                right = mid - 1

        min_ = arr[right]
        overflow = k - (right * arr[right] - preSum[right])  # 填平后多出来的
        div, mod = 0, 0
        if right:
            div, mod = divmod(overflow, right)
        min_ += div

        res = 1
        for i, num in enumerate(sorted(nums)):
            res *= max(num, min_ + int(i < mod))
            res %= MOD
        return res


print(Solution().maximumProduct(nums=[6, 3, 3, 2], k=2))
