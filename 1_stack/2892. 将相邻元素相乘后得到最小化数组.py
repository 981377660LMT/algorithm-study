# 2892. 将相邻元素相乘后得到最小化数组
# https://leetcode.cn/problems/minimizing-array-after-replacing-pairs-with-their-product/description/
# 给定一个整数数组 nums 和一个整数 k，你可以对数组执行以下操作任意次数：


# 选择数组中的两个 相邻 元素，例如 x 和 y，使得 x * y <= k ，并用一个值为 x * y 的 单个元素 替换它们（
# 例如，在一次操作中，数组 [1, 2, 2, 3]，其中 k = 5 可以变为 [1, 4, 3] 或 [2, 2, 3]，但不能变为 [1, 2, 6]）。
# 返回 经过任意次数的操作后， nums 的 最小 可能长度。
# 1 <= nums.length <= 1e5
# 0 <= nums[i] <= 1e9
# 1 <= k <= 1e9

# !如果栈顶两个元素可以合并，就必须要合并


from typing import List


class Solution:
    def minArrayLength(self, nums: List[int], k: int) -> int:
        stack = []
        for num in nums:
            stack.append(num)
            while len(stack) >= 2 and stack[-1] * stack[-2] <= k:
                a = stack.pop()
                b = stack.pop()
                stack.append(a * b)
        return len(stack)
