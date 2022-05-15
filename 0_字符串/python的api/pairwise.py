from itertools import pairwise
from typing import List


class Solution:
    def maxConsecutive(self, bottom: int, top: int, special: List[int]) -> int:
        """返回不含特殊楼层的 最大 连续楼层数
        
        分割区间
        """
        special = sorted(special + [bottom - 1, top + 1])
        return max(b - a - 1 for a, b in pairwise(special))
        return max(b - a - 1 for a, b in zip(special, special[1:]))


# pairwise(iterator)
#  等价于
# zip(iterator, iterator[1:])
