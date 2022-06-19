from typing import List

# 给你一个正整数数组 nums，请你移除 最短 子数组（可以为 空），使得剩余元素的 和 能被 p 整除。 不允许 将整个数组都移除。
# 请你返回你需要移除的最短子数组的长度，如果无法满足题目要求，返回 -1 。

# 1 <= nums.length <= 105

# 思路：
# 1.寻找最短的一段子数组和S S与sum(nums) 模p同余
# 2.和按照模p分类，记录`最近`的前缀和


class Solution:
    def minSubarray(self, nums: List[int], p: int) -> int:
        need = sum(nums) % p
        if need == 0:
            return 0

        prefix = {0: -1}
        curSum = 0
        res = n = len(nums)
        for i, num in enumerate(nums):
            curSum = (curSum + num) % p
            if (curSum - need) % p in prefix:
                res = min(res, i - prefix[(curSum - need) % p])
            prefix[curSum] = i
        return res if res < n else -1


print(Solution().minSubarray(nums=[3, 1, 4, 2], p=6))
# 输出：1
# 解释：nums 中元素和为 10，不能被 p 整除。我们可以移除子数组 [4] ，剩余元素的和为 6 。
