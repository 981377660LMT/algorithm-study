from itertools import pairwise
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 如果数组的每一对相邻元素都是两个奇偶性不同的数字，则该数组被认为是一个 特殊数组 。

# 周洋哥有一个整数数组 nums 和一个二维整数矩阵 queries，对于 queries[i] = [fromi, toi]，请你帮助周洋哥检查子数组 nums[fromi..toi] 是不是一个 特殊数组 。


# 返回布尔数组 answer，如果 nums[fromi..toi] 是特殊数组，则 answer[i] 为 true ，否则，answer[i] 为 false 。


# [0,1] => 0
# [1,3] => sum(1,2)
class Solution:
    def isArraySpecial(self, nums: List[int], queries: List[List[int]]) -> List[bool]:
        nums = [v & 1 for v in nums]
        preSum = [0]
        for a, b in pairwise(nums):
            preSum.append(preSum[-1] + int(a == b))
        res = []
        for a, b in queries:
            ok = preSum[b] == preSum[a]
            res.append(ok)
        return res


# nums = [4,3,1,6], queries = [[0,2],[2,3]]
print(Solution().isArraySpecial([4, 3, 1, 6], [[0, 2], [2, 3]]))  # [false,true]
