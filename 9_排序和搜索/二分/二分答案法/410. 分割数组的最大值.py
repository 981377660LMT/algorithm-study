# 410. 分割数组的最大值
# 给定一个非负整数数组 nums 和一个整数 k ，你需要将这个数组分成 k 个非空的连续子数组。
# !设计一个算法使得这 k 个子数组各自和的最大值最小。
# n<=1000 0<=nums[i]<=1e6


from typing import List


class Solution:
    def splitArray(self, nums: List[int], k: int) -> int:
        def check(mid: int) -> bool:
            """每段不超过mid,是否能分成<=k段"""
            count = 1
            curSum = 0
            for num in nums:
                if num > mid:
                    return False
                curSum += num
                if curSum > mid:
                    count += 1
                    curSum = num
            return count <= k

        left, right = 0, len(nums) * max(nums)
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1
        return left


assert Solution().splitArray([7, 2, 5, 10, 8], 2) == 18
