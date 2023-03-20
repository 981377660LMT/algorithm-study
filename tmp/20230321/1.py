from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个 正 整数 n 。

# 用 even 表示在 n 的二进制形式（下标从 0 开始）中值为 1 的偶数下标的个数。

# 用 odd 表示在 n 的二进制形式（下标从 0 开始）中值为 1 的奇数下标的个数。

# 返回整数数组 answer ，其中 answer = [even, odd] 。
class Solution:
    def evenOddBit(self, n: int) -> List[int]:
        bin_ = bin(int(n))[2:]
        even, odd = 0, 0
        for i, c in enumerate(bin_[::-1]):
            if c == "1":
                if i % 2 == 0:
                    even += 1
                else:
                    odd += 1
        return [even, odd]
