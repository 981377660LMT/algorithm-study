from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你两个 正整数 n 和 k。有 n 个编号从 0 到 n - 1 的孩子按顺序从左到右站成一队。

# 最初，编号为 0 的孩子拿着一个球，并且向右传球。每过一秒，拿着球的孩子就会将球传给他旁边的孩子。一旦球到达队列的 任一端 ，即编号为 0 的孩子或编号为 n - 1 的孩子处，传球方向就会 反转 。


# 返回 k 秒后接到球的孩子的编号。
class Solution:
    def numberOfChild(self, n: int, k: int) -> int:
        dir = 1
        cur = 0
        for _ in range(k):
            cur += dir
            if cur == n - 1 or cur == 0:
                dir = -dir
        return cur
