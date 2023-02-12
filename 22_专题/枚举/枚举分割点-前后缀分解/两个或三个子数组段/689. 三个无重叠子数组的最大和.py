from itertools import accumulate
from typing import List

# 给定数组 nums 由`正整数`组成,找到三个互不重叠的子数组的最大和。
# 每个子数组的长度为k，我们要使这3*k个项的和最大化。
# nums.length的范围在[1, 20000]之间。
# 返回每个区间起始索引的列表（索引从 0 开始）。如果有多个结果，返回字典序最小的一个。


class Solution:
    # 求和
    def maxSumOfThreeSubarrays1(self, nums: List[int], k: int) -> int:
        n = len(nums)
        preSum = [0] + list(accumulate(nums))
        kSum = [preSum[i] - preSum[i - k] for i in range(k, n + 1)]
        preMax = list(accumulate(kSum, max))
        suffixMax = list(accumulate(kSum[::-1], max))[::-1]
        return max(preMax[i - k] + kSum[i] + suffixMax[i + k] for i in range(k, len(kSum) - k))

    # 求对应子数组(起点)
    def maxSumOfThreeSubarrays(self, nums: List[int], k: int) -> List[int]:
        n = len(nums)
        preSum = [0] + list(accumulate(nums))
        kSum = [preSum[i] - preSum[i - k] for i in range(k, n + 1)]

        preMax = list(accumulate(kSum, max))
        preMaxIndex = [0] * len(preMax)
        for i in range(1, len(preMax)):
            if preMax[i] > preMax[i - 1]:
                preMaxIndex[i] = i
            else:
                preMaxIndex[i] = preMaxIndex[i - 1]

        suffixMax = list(accumulate(kSum[::-1], max))[::-1]
        suffixMaxIndex = [n - k] * len(preMax)
        for i in range(n - k - 1, -1, -1):
            # 字典序最小
            if suffixMax[i] > suffixMax[i + 1] or (
                suffixMax[i] == suffixMax[i + 1] and kSum[i] == kSum[suffixMaxIndex[i + 1]]
            ):
                suffixMaxIndex[i] = i
            else:
                suffixMaxIndex[i] = suffixMaxIndex[i + 1]

        res, maxSum = [-1, -1, -1], -int(1e20)
        for i in range(k, len(kSum) - k):
            if preMax[i - k] + kSum[i] + suffixMax[i + k] > maxSum:
                res = [preMaxIndex[i - k], i, suffixMaxIndex[i + k]]
                maxSum = preMax[i - k] + kSum[i] + suffixMax[i + k]

        return res


# print(Solution().maxSumOfThreeSubarrays([1, 2, 1, 2, 6, 7, 5, 1], 2))
# # 输出: [0, 3, 5]
# # 解释: 子数组 [1, 2], [2, 6], [7, 5] 对应的起始索引为 [0, 3, 5]。
# # 我们也可以取 [2, 1], 但是结果 [1, 3, 5] 在字典序上更大。
# print(Solution().maxSumOfThreeSubarrays([1, 2, 1, 2, 1, 2, 1, 2, 1], 2))
# print(Solution().maxSumOfThreeSubarrays([4, 5, 10, 6, 11, 17, 4, 11, 1, 3], 1))
print(
    Solution().maxSumOfThreeSubarrays(
        [17, 9, 3, 2, 7, 10, 20, 1, 13, 4, 5, 16, 4, 1, 17, 6, 4, 19, 8, 3], 4
    )
)
