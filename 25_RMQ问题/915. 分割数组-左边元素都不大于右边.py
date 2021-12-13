from typing import List

# 给定一个数组 A，将其划分为两个连续子数组 left 和 right
# left 中的每个元素都小于或等于 right 中的每个元素。
# left 和 right 都是非空的。
# left 的长度要尽可能小。
# 在完成这样的分组后返回 left 的长度


# 不检验 all(L <= R for L in left for R in right)，而是检验 max(left) <= min(right)
# 辅助数组：扫两遍，记录leftMax和rightMin  => 类似接雨水那题的数组
class Solution:
    def partitionDisjoint(self, nums: List[int]) -> int:
        n = len(nums)
        leftMax = [nums[0]] * n
        rightMin = [nums[-1]] * n

        curMax = nums[0]
        for i in range(n):
            curMax = max(curMax, nums[i])
            leftMax[i] = curMax

        curMin = nums[-1]
        for i in range(n - 1, -1, -1):
            curMin = min(curMin, nums[i])
            rightMin[i] = curMin

        for i in range(1, n):
            if leftMax[i - 1] <= rightMin[i]:
                return i


print(Solution().partitionDisjoint([1, 1, 1, 0, 6, 12]))
