from typing import List


class Solution:
    def maxSumTwoNoOverlap(self, nums: List[int], left: int, right: int):
        # L在i左边，M在i右边或L在i右边，M在i左边。
        # 对于每个结尾的i，按照两种情况计算最大的L和M和，并更新答案
        n = len(nums)
        preSum = [0] * (n + 1)
        for i in range(n):
            preSum[i + 1] = preSum[i] + nums[i]

        res = preSum[left + right]

        # 维护取前面的最大值
        lMax = preSum[left]
        rMax = preSum[right]
        # i代表当前位于右边的数组的末尾索引
        for i in range(left + right, n + 1):
            # 当L在M前时，i代表M的最后一个索引,此时M已确定
            lMax = max(lMax, preSum[i - right] - preSum[i - right - left])
            cand1 = lMax + (preSum[i] - preSum[i - right])
            # 当L在M后时，i代表L的最后一个索引，此时L已确定
            rMax = max(rMax, preSum[i - left] - preSum[i - left - right])
            cand2 = rMax + (preSum[i] - preSum[i - left])
            res = max(res, cand1, cand2)
        return res


print(Solution().maxSumTwoNoOverlap([0, 6, 5, 2, 2, 5, 1, 9, 4], 1, 2))
print(Solution().maxSumTwoNoOverlap([1, 0, 1], 1, 1))
print(Solution().maxSumTwoNoOverlap([1, 0, 3], 1, 2))
