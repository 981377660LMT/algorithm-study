from typing import List
from itertools import accumulate


# 请你在 arr 中找 两个互不重叠的子数组 且它们的和都等于 target 。
# 可能会有多种方案，请你返回满足要求的两个子数组长度和的 最小值 。
# 请返回满足要求的最小长度和，如果无法找到这样的两个子数组，请返回 -1 。

# 1 <= arr.length <= 10^5
# 1 <= arr[i] <= 1000  正数 前缀和递增

# 思路：哈希表+前缀和
# 先遍历一遍记录前缀和，然后扫一遍对每个位置，找到 curSum-target 与 curSum+target 对应的索引
INF = 0x3F3F3F3F


class Solution:
    def minSumOfLengths(self, arr: List[int], target: int) -> int:
        preSum, curSum = dict({0: -1}), 0
        for i, num in enumerate(arr):
            curSum += num
            preSum[curSum] = i

        curSum = 0
        leftSize = INF
        res = INF
        for i, num in enumerate(arr):
            curSum += num
            leftSum, rightSum = curSum - target, curSum + target
            if leftSum in preSum:
                leftSize = min(leftSize, i - preSum[leftSum])
            if rightSum in preSum:
                res = min(res, preSum[rightSum] - i + leftSize)

        return res if res != INF else -1


print(Solution().minSumOfLengths([3, 4, 7, 7], 7))
print(Solution().minSumOfLengths(arr=[3, 2, 2, 4, 3], target=3))
# 不一定先看就最短
print(Solution().minSumOfLengths(arr=[2, 1, 3, 3, 2, 3, 1], target=6))

