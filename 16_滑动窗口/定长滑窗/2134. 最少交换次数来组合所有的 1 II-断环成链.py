from typing import List


class Solution:
    def minSwaps(self, nums: List[int]) -> int:
        winLen = nums.count(1)
        nums *= 2  # 破环成链
        res = int(1e20)
        one = 0
        for i, cur in enumerate(nums):
            if cur == 1:
                one += 1
            if i >= winLen:
                if nums[i - winLen] == 1:
                    one -= 1
            if i >= winLen - 1:
                res = min(res, winLen - one)
        return res


print(Solution().minSwaps(nums=[0, 1, 0, 1, 1, 0, 0]))
print(Solution().minSwaps(nums=[0, 1, 1, 1, 0, 0, 1, 1, 0]))
print(Solution().minSwaps(nums=[1, 1, 0, 0, 1]))
