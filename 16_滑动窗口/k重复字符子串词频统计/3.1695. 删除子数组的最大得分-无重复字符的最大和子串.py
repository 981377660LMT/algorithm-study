from typing import List


class Solution:
    def maximumUniqueSubarray(self, nums: List[int]) -> int:
        res = 0
        curSum = 0
        left, right = 0, 0
        visited = set()
        while right < len(nums):
            if nums[right] not in visited:
                curSum += nums[right]
                visited.add(nums[right])
                right += 1
                res = max(res, curSum)
            else:
                visited.discard(nums[left])
                curSum -= nums[left]
                left += 1

        return res


print(Solution().maximumUniqueSubarray(nums=[4, 2, 4, 5, 6]))
# 输出：17
# 解释：最优子数组是 [2,4,5,6]
