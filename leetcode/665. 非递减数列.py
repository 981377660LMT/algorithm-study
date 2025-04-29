# 665. 非递减数列
# 给你一个长度为 n 的整数数组 nums ，请你判断在 最多 改变 1 个元素的情况下，该数组能否变成一个非递减数列。

from typing import List


class Solution:
    def checkPossibility(self, nums: List[int]) -> bool:
        """
        贪心扫描，记录逆序对次数 count。遇到 nums[i] < nums[i-1] 时：
          - 如果已经有一次修改，直接返回 False。
          - 否则尝试“修正”一次：
            • 如果 i<2 或 nums[i] >= nums[i-2]，则把 nums[i-1] 调低到 nums[i]（等价于修改前一个元素）。
            • 否则把 nums[i] 调高到 nums[i-1]（修改当前元素）。
        最终若 count ≤ 1 则可行。
        时间 O(n)，空间 O(1)。
        """
        count = 0
        for i in range(1, len(nums)):
            if nums[i] < nums[i - 1]:
                count += 1
                if count > 1:
                    return False
                # 尝试修改 nums[i-1] 或 nums[i]
                if i < 2 or nums[i] >= nums[i - 2]:
                    # 安全地把前一个元素调低
                    nums[i - 1] = nums[i]
                else:
                    # 否则只能把当前元素调高
                    nums[i] = nums[i - 1]
        return True
