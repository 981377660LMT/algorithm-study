from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一个下标从 0 开始、由正整数组成的数组 nums 。

# 你可以在数组上执行下述操作 任意 次：

# 选中一个同时满足 0 <= i < nums.length - 1 和 nums[i] <= nums[i + 1] 的整数 i 。将元素 nums[i + 1] 替换为 nums[i] + nums[i + 1] ，并从数组中删除元素 nums[i] 。
# 返回你可以从最终数组中获得的 最大 元素的值。


class Solution:
    def maxArrayValue(self, nums: List[int]) -> int:
        for i in range(len(nums) - 2, -1, -1):
            if nums[i] <= nums[i + 1]:
                nums[i] = nums[i] + nums[i + 1]
                nums[i + 1] = 0
        return max(nums)
