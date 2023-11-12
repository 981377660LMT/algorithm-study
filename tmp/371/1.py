from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的整数数组 nums 。如果一对整数 x 和 y 满足以下条件，则称其为 强数对 ：

# |x - y| <= min(x, y)
# 你需要从 nums 中选出两个整数，且满足：这两个整数可以形成一个强数对，并且它们的按位异或（XOR）值是在该数组所有强数对中的 最大值 。

# 返回数组 nums 所有可能的强数对中的 最大 异或值。


# 注意，你可以选择同一个整数两次来形成一个强数对。


class Solution:
    def maximumStrongPairXor(self, nums: List[int]) -> int:
        res = 0
        for a in nums:
            for b in nums:
                if abs(a - b) <= min(a, b):
                    res = max(res, a ^ b)
        return res
