from typing import List
from itertools import accumulate

# 1 <= n <= 105
# 构建并返回一个长度为 n 的数组 avgs ，其中 avgs[i] 是以下标 i 为中心的子数组的 半径为 k 的子数组平均值 。

# 1. 滑动窗口
# 2. 前缀和
class Solution:
    def getAverages1(self, nums: List[int], k: int) -> List[int]:
        res = [-1] * len(nums)
        curSum = 0  # range sum
        for i, x in enumerate(nums):
            curSum += x
            if i >= 2 * k + 1:
                curSum -= nums[i - (2 * k + 1)]
            if i + 1 >= 2 * k + 1:
                res[i - k] = curSum // (2 * k + 1)
        return res

    def getAverages(self, nums: List[int], k: int) -> List[int]:
        preSum = [0] + list(accumulate(nums))
        res = [-1] * len(nums)
        for i, _ in enumerate(nums):
            if k <= i < len(nums) - k:
                res[i] = (preSum[i + k + 1] - preSum[i - k]) // (2 * k + 1)
        return res


print(Solution().getAverages(nums=[7, 4, 3, 9, 1, 8, 5, 2, 6], k=3))
# 输出：[-1,-1,-1,5,4,4,-1,-1,-1]
# 解释：
# - avg[0]、avg[1] 和 avg[2] 是 -1 ，因为在这几个下标前的元素数量都不足 k 个。
# - 中心为下标 3 且半径为 3 的子数组的元素总和是：7 + 4 + 3 + 9 + 1 + 8 + 5 = 37 。
#   使用截断式 整数除法，avg[3] = 37 / 7 = 5 。
# - 中心为下标 4 的子数组，avg[4] = (4 + 3 + 9 + 1 + 8 + 5 + 2) / 7 = 4 。
# - 中心为下标 5 的子数组，avg[5] = (3 + 9 + 1 + 8 + 5 + 2 + 6) / 7 = 4 。
# - avg[6]、avg[7] 和 avg[8] 是 -1 ，因为在这几个下标后的元素数量都不足 k 个。

