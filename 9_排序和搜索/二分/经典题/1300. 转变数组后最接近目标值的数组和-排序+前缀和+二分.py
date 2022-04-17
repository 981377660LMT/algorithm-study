from typing import List
from itertools import accumulate
from bisect import bisect_left

# 请你返回一个整数 value ，使得将数组中所有大于 value 的值变成 value 后，数组的和最接近  target
# 一眼二分 1 <= arr.length <= 10^4

# 如果有多种使得和最接近 target 的方案，请你返回这些整数中的最小值。
class Solution:
    def findBestValue(self, arr: List[int], target: int) -> int:
        def calSum(mid: int) -> int:
            index = bisect_left(arr, mid)  # 大于mid的第一个数
            return preSum[index] + (n - index) * mid

        n = len(arr)
        arr = sorted(arr)
        preSum = [0] + list(accumulate(arr))

        left, right = 0, max(arr)
        while left <= right:
            mid = (left + right) >> 1
            if calSum(mid) >= target:
                right = mid - 1
            else:
                left = mid + 1

        cand1, cand2 = left - 1, left
        diff1, diff2 = abs(target - calSum(cand1)), abs(target - calSum(cand2))
        return cand1 if diff1 <= diff2 else cand2


print(Solution().findBestValue(arr=[4, 9, 3], target=10))
# 输出：3
# 解释：当选择 value 为 3 时，数组会变成 [3, 3, 3]，和为 9 ，这是最接近 target 的方案。
