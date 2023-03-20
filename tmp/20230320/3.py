from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个数组 nums ，它包含若干正整数。

# 一开始分数 score = 0 ，请你按照下面算法求出最后分数：

# 从数组中选择最小且没有被标记的整数。如果有相等元素，选择下标最小的一个。
# 将选中的整数加到 score 中。
# 标记 被选中元素，如果有相邻元素，则同时标记 与它相邻的两个元素 。
# 重复此过程直到数组中所有元素都被标记。
# 请你返回执行上述算法后最后的分数。
class Solution:
    def findScore(self, nums: List[int]) -> int:
        sl = SortedList((v, i) for i, v in enumerate(nums))
        dead = set()
        res = 0
        while sl:
            v, i = sl.pop(0)
            res += v
            dead.add(i)

            if i > 0 and i - 1 not in dead:
                sl.remove((nums[i - 1], i - 1))
                dead.add(i - 1)
            if i < len(nums) - 1 and i + 1 not in dead:
                sl.remove((nums[i + 1], i + 1))
                dead.add(i + 1)
        return res


# nums = [2,1,3,4,5,2]
print(Solution().findScore([2, 1, 3, 4, 5, 2]))
