from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个正整数数组 nums 和一个整数 k 。

# 一次操作中，你可以将数组的最后一个元素删除，将该元素添加到一个集合中。


# 请你返回收集元素 1, 2, ..., k 需要的 最少操作次数 。
class Solution:
    def minOperations(self, nums: List[int], k: int) -> int:
        res = set()
        target = set(range(1, k + 1))

        for i in range(len(nums) - 1, -1, -1):
            res.add(nums[i])
            if res & target == target:
                return len(nums) - i


# nums = [3,1,5,4,2], k = 2
print(Solution().minOperations([3, 1, 5, 4, 2], 2))
