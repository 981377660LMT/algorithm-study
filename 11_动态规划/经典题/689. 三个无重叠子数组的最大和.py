from typing import List

# 给定数组 nums 由`正整数`组成,找到三个互不重叠的子数组的最大和。
# 每个子数组的长度为k，我们要使这3*k个项的和最大化。
# nums.length的范围在[1, 20000]之间。
# 返回每个区间起始索引的列表（索引从 0 开始）。如果有多个结果，返回字典序最小的一个。


# 暴力滑窗
# 向右滑动的同时严格大于(>)避免了重复结果
class Solution:
    def maxSumOfThreeSubarrays(self, nums: List[int], k: int) -> List[int]:

        # Best single, double, and triple sequence found so far
        best1Win = 0
        best2Win = [0, k]
        best3Win = [0, k, k * 2]

        # Sums of each window
        win1Sum = sum(nums[0:k])
        win2Sum = sum(nums[k : k * 2])
        win3Sum = sum(nums[k * 2 : k * 3])

        # Sums of combined best windows
        best1WinSum = win1Sum
        best2WinSum = win1Sum + win2Sum
        best3WinSum = win1Sum + win2Sum + win3Sum

        # Current window positions
        i1 = 1
        i2 = k + 1
        i3 = k * 2 + 1
        while i3 <= len(nums) - k:
            # Update the three sliding windows
            win1Sum = win1Sum - nums[i1 - 1] + nums[i1 + k - 1]
            win2Sum = win2Sum - nums[i2 - 1] + nums[i2 + k - 1]
            win3Sum = win3Sum - nums[i3 - 1] + nums[i3 + k - 1]

            # Update best single window
            if win1Sum > best1WinSum:
                best1Win = i1
                best1WinSum = win1Sum

            # Update best two windows
            if win2Sum + best1WinSum > best2WinSum:
                best2Win = [best1Win, i2]
                best2WinSum = win2Sum + best1WinSum

            # Update best three windows
            if win3Sum + best2WinSum > best3WinSum:
                best3Win = best2Win + [i3]
                best3WinSum = win3Sum + best2WinSum

            # Update the current positions
            i1 += 1
            i2 += 1
            i3 += 1

        return best3Win


print(Solution().maxSumOfThreeSubarrays([1, 2, 1, 2, 6, 7, 5, 1], 2))
# 输出: [0, 3, 5]
# 解释: 子数组 [1, 2], [2, 6], [7, 5] 对应的起始索引为 [0, 3, 5]。
# 我们也可以取 [2, 1], 但是结果 [1, 3, 5] 在字典序上更大。
