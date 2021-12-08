from typing import List


# 给你一个 非递减 有序整数数组 nums 。(暗示前缀+二分)
# 请你建立并返回一个整数数组 result，它跟 nums 长度相同，且result[i] 等于 nums[i] 与数组中所有其他元素差的绝对值之和。


# i位置为基准，左右绝对值展开计算
# postsum[i] - presum[i] - nums[i] * (n - 2 * i - 1);
class Solution:
    def getSumAbsoluteDifferences(self, nums: List[int]) -> List[int]:
        preSum = [0]
        for num in nums:
            preSum.append(preSum[-1] + num)
        n = len(nums)
        return [
            (num * (i + 1) - preSum[i + 1]) + (preSum[-1] - preSum[i] - (n - i) * num)
            for i, num in enumerate(nums)
        ]


print(Solution().getSumAbsoluteDifferences(nums=[2, 3, 5]))
# 输出：[4,3,5]
# 解释：假设数组下标从 0 开始，那么
# result[0] = |2-2| + |2-3| + |2-5| = 0 + 1 + 3 = 4，
# result[1] = |3-2| + |3-3| + |3-5| = 1 + 0 + 2 = 3，
# result[2] = |5-2| + |5-3| + |5-5| = 3 + 2 + 0 = 5。
