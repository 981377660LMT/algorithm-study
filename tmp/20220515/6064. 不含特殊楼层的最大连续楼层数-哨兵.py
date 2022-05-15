from itertools import pairwise
from typing import List, Tuple


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def maxConsecutive(self, bottom: int, top: int, special: List[int]) -> int:
        """返回不含特殊楼层的 最大 连续楼层数
        
        分割区间-记录pre，添加哨兵
        """
        special.append(top + 1)
        special = sorted(special)
        pre = bottom
        res = 0
        for bad in special:
            cur = (bad - 1) - pre + 1
            res = max(res, cur)
            pre = bad + 1
        return res

    def maxConsecutive2(self, bottom: int, top: int, special: List[int]) -> int:
        """返回不含特殊楼层的 最大 连续楼层数
        
        分割区间
        """
        special = sorted(special + [bottom - 1, top + 1])
        return max(b - a - 1 for a, b in pairwise(special))
        return max(b - a - 1 for a, b in zip(special, special[1:]))


# 输入：bottom = 2, top = 9, special = [4,6]
# 输出：3
# 解释：下面列出的是不含特殊楼层的连续楼层范围：
# - (2, 3) ，楼层数为 2 。
# - (5, 5) ，楼层数为 1 。
# - (7, 9) ，楼层数为 3 。
# 因此，返回最大连续楼层数 3 。

