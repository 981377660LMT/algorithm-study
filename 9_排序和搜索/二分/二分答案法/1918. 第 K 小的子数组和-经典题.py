from typing import List


# 这也能二分?:看到K-th元素的题目第一反应一般都是二分搜索，此时我们可以通过二分搜索查询是否是第K个答案。
class Solution:
    def kthSmallestSubarraySum(self, nums: List[int], k: int) -> int:
        def countNGT(mid) -> int:
            """"和小于等于mid的子数组数"""

            res, curSum, left = 0, 0, 0
            for right in range(len(nums)):
                curSum += nums[right]
                while left < len(nums) and curSum > mid:
                    curSum -= nums[left]
                    left += 1
                res += right - left + 1
            return res

        left, right = min(nums), sum(nums)
        while left <= right:
            mid = (left + right) >> 1
            # 找最左，尽量把右边移过来
            if countNGT(mid) < k:
                left = mid + 1
            else:
                right = mid - 1
        return left


print(Solution().kthSmallestSubarraySum(nums=[2, 1, 3], k=4))
# 输出: 3
# 解释: [2,1,3] 的子数组为：
# - [2] 和为 2
# - [1] 和为 1
# - [3] 和为 3
# - [2,1] 和为 3
# - [1,3] 和为 4
# - [2,1,3] 和为 6
# 最小子数组和的升序排序为 1, 2, 3, 3, 4, 6。 第 4 小的子数组和为 3 。
