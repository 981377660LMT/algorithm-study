from itertools import groupby
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个 二进制字符串 s。

# 你可以对这个字符串执行 任意次 下述操作：
        c1 = 0
        res = 0
        for t, g in groupby(s):
            ln = len(list(g))
            if t == '1':
                c1 += ln
            else:
                res += c1
        return res

# 选择字符串中的任一下标 i（ i + 1 < s.length ），该下标满足 s[i] == '1' 且 s[i + 1] == '0'。
# 将字符 s[i] 向 右移 直到它到达字符串的末端或另一个 '1'。例如，对于 s = "010010"，如果我们选择 i = 1，结果字符串将会是 s = "000110"。
# 返回你能执行的 最大 操作次数。
class Solution:
    def maxOperations(self, s: str) -> int:
        groups = [(char, len(list(group))) for char, group in groupby(s)]
        onesGroupLen = [groupLen for char, groupLen in groups if char == "1"]
        curSum = 0
        res = 0
        for v in onesGroupLen:
            res += curSum
            curSum += v
        if s[-1] == "0":
            res += curSum
        return res


#  s = "1001101"

print(Solution().maxOperations("1001101"))  # 2
# "110"
print(Solution().maxOperations("110"))  # 2
