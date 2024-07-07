from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个正整数 n。

# 如果一个二进制字符串 x 的所有长度为 2 的子字符串中包含 至少 一个 "1"，则称 x 是一个 有效 字符串。

# 返回所有长度为 n 的 有效 字符串，可以以任意顺序排列。


class Solution:
    def validStrings(self, n: int) -> List[str]:
        if n == 1:
            return ["0", "1"]
        res = []
        for state in range(1, 1 << n):
            ok = [False] * n
            for i in range(n):
                if (state >> i) & 1:
                    ok[i] = True
            if all(ok[i] or ok[i + 1] for i in range(n - 1)):
                res.append("".join("1" if ok[i] else "0" for i in range(n)))
        return res
