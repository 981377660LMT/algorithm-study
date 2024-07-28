from itertools import groupby
from tkinter import N
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个二进制字符串 s。

# 请你统计并返回其中 1 显著 的 子字符串 的数量。


# 如果字符串中 1 的数量 大于或等于 0 的数量的 平方，则认为该字符串是一个 1 显著 的字符串


# !0不多，枚举0的个数?


def max2(a: int, b: int) -> int:
    return a if a > b else b


def getNextIndex(arr: List[int]) -> List[int]:
    """
    获取数组中相同元素的后一个元素的位置.不存在则返回-1.
    """
    pool = dict()
    for i, v in enumerate(arr):
        arr[i] = pool.setdefault(v, len(pool))

    n = len(arr)
    nexts, valueNexts = [-1] * n, [-1] * len(pool)
    for i in range(n - 1, -1, -1):
        v = arr[i]
        if valueNexts[v] != -1:
            nexts[i] = valueNexts[v]
        valueNexts[v] = i
    return nexts


class Solution:
    def numberOfSubstrings(self, s: str) -> int:
        n = len(s)
        nums = [1 if ch == "1" else 0 for ch in s]
        nexts = getNextIndex(nums)
        print(nexts)


# s = "00011"

print(Solution().numberOfSubstrings(s="101101"))  # 3
