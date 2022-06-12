from typing import List


# !注意子数组和前缀和求出 滑窗更新
class Solution:
    def countSubarrays(self, nums: List[int], k: int) -> int:
        """根据单调性 枚举子数组右端点，去看对应的合法左端点的个数"""
        curSum, left, res = 0, 0, 0
        for right, num in enumerate(nums):
            curSum += num
            while curSum * (right - left + 1) >= k:
                curSum -= nums[left]
                left += 1
            res += right - left + 1
        return res


print(Solution().countSubarrays(nums=[2, 1, 4, 3, 5], k=10))
print(Solution().countSubarrays(nums=[1, 1, 1], k=5))
