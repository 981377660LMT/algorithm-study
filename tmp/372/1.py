from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你三个字符串 s1、s2 和 s3。 你可以根据需要对这三个字符串执行以下操作 任意次数 。

# 在每次操作中，你可以选择其中一个长度至少为 2 的字符串 并删除其 最右位置上 的字符。


# 如果存在某种方法能够使这三个字符串相等，请返回使它们相等所需的 最小 操作次数；否则，返回 -1。
class Solution:
    def findMinimumOperations(self, s1: str, s2: str, s3: str) -> int:
        if len(set([s1[0], s2[0], s3[0]])) != 1:
            return -1
        n1, n2, n3 = len(s1), len(s2), len(s3)
        ptr = 0
        while ptr < n1 and ptr < n2 and ptr < n3 and s1[ptr] == s2[ptr] == s3[ptr]:
            ptr += 1
        return n1 - ptr + n2 - ptr + n3 - ptr
