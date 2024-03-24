from math import ceil
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个正整数 k 。最初，你有一个数组 nums = [1] 。

# 你可以对数组执行以下 任意 操作 任意 次数（可能为零）：


# 选择数组中的任何一个元素，然后将它的值 增加 1 。
# 复制数组中的任何一个元素，然后将它附加到数组的末尾。
# 返回使得最终数组元素之 和 大于或等于 k 所需的 最少 操作次数。
class Solution:
    def minOperations(self, k: int) -> int:
        if k == 1:
            return 0
        # 枚举增加的次数
        res = k
        for add in range(k + 1):
            base = 1 + add
            div = ceil((k - base) / base)
            res = min(res, add + div)
        return res
