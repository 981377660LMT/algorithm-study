from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的整数数组 nums 和一个整数 k 。

# nums 中的 K-or 是一个满足以下条件的非负整数：

# 只有在 nums 中，至少存在 k 个元素的第 i 位值为 1 ，那么 K-or 中的第 i 位的值才是 1 。
# 返回 nums 的 K-or 值。


# 注意 ：对于整数 x ，如果 (2i AND x) == 2i ，则 x 中的第 i 位值为 1 ，其中 AND 为按位与运算符。
class Solution:
    def findKOr(self, nums: List[int], k: int) -> int:
        bitCounter = [0] * 32
        for num in nums:
            for i in range(32):
                if num & (1 << i):
                    bitCounter[i] += 1
        return sum(1 << i for i in range(32) if bitCounter[i] >= k)
