from typing import List

INF = int(1e20)


class Solution:
    def solve(self, nums: List[int], target: int) -> int:
        """nums 正整数数组"""
        """头尾pop数 问是否能正好删除target"""
        """等价于滑窗内和能凑成sum(nums)-target时的最长滑窗"""
        if target == 0:
            return 0

        target = sum(nums) - target
        left, curSum, res = 0, 0, -INF
        if target == 0:
            res = 0

        for right, num in enumerate(nums):
            curSum += num
            if curSum == target:
                res = max(res, right - left + 1)
            elif curSum > target:
                while left < right and curSum > target:
                    curSum -= nums[left]
                    left += 1
                if curSum == target:
                    res = max(res, right - left + 1)

        return -1 if res == -INF else len(nums) - res


print(Solution().solve([1], 1))
print(Solution().solve([1, 2], 1))
