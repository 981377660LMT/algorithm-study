from typing import List

# 找出该数组中满足其和 ≥ target 的长度最小的 连续子数组
class Solution:
    # nlogn解法 二分答案
    def minSubArrayLen(self, s: int, nums: List[int]) -> int:
        left, right, res = 0, len(nums), 0

        def check(size):
            sum_size = 0
            for i in range(len(nums)):
                sum_size += nums[i]
                if i >= size:
                    sum_size -= nums[i - size]
                if sum_size >= s:
                    return True
            return False

        while left <= right:
            mid = (left + right) // 2  # 滑动窗口大小
            if check(mid):  # 如果这个大小的窗口可以那么就缩小
                res = mid
                right = mid - 1
            else:  # 否则就增大窗口
                left = mid + 1
        return res
