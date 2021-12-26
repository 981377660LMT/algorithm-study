from typing import List

# 非递减 有序整数数组 nums
# result[i] 等于 sum(|nums[i]-nums[j]|)


# 前缀和，以当前数为分割点，左边的数都比当前数小，右边的数都比当前数大。
# postsum[i] - presum[i] - nums[i] * (n - 2 * i - 1);
class Solution:
    def getSumAbsoluteDifferences(self, nums: List[int]) -> List[int]:
        n = len(nums)
        preSum = [0]
        for num in nums:
            preSum.append(preSum[-1] + num)
        return [
            (num * (i + 1) - preSum[i + 1]) + (preSum[n] - preSum[i] - (n - i) * num)
            for i, num in enumerate(nums)
        ]


print(Solution().getSumAbsoluteDifferences(nums=[2, 3, 5]))
print(Solution().getSumAbsoluteDifferences(nums=[2]))
# 输出：[4,3,5]
# 解释：假设数组下标从 0 开始，那么
# result[0] = |2-2| + |2-3| + |2-5| = 0 + 1 + 3 = 4，
# result[1] = |3-2| + |3-3| + |3-5| = 1 + 0 + 2 = 3，
# result[2] = |5-2| + |5-3| + |5-5| = 3 + 2 + 0 = 5。
