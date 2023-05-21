from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个正整数 n ，请你返回 n 的 惩罚数 。

# n 的 惩罚数 定义为所有满足以下条件 i 的数的平方和：


# 1 <= i <= n
# i * i 的十进制表示的字符串可以分割成若干连续子字符串，且这些子字符串对应的整数值之和等于 i 。


@lru_cache(None)
def check(num: int) -> bool:
    n2 = str(num * num)
    # 枚举分割点
    n = len(n2)
    for state in range(1 << (n - 1)):
        # 枚举分割方案
        curSplit = []
        for i in range(n - 1):
            if state & (1 << i):
                curSplit.append(i + 1)
        curSplit = [0] + curSplit + [n]
        curSum = sum(int(n2[curSplit[i] : curSplit[i + 1]]) for i in range(len(curSplit) - 1))
        if curSum == num:
            return True
    return False


class Solution:
    def punishmentNumber(self, n: int) -> int:
        return sum(i * i for i in range(1, n + 1) if check(i))
