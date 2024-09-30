from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

from typing import List


def makeDp(seq: List[int], rev=True) -> List[int]:
    n = len(seq)
    dp = [0] * (n + 1)
    for i in range(n):
        cur = seq[i]
        # your code here
    return dp


pre, suf = makeDp(nums, False), makeDp(nums[::-1])[::-1]
res = 0
# 前后缀分解容易错，完善模板及相关题
for i in range(len(pre)):
    res += pre[i] * suf[i]  # [0,i) [i,n)


# 失配一次的子序列
# !子序列+前后缀，以前好像出过???
class Solution:
    def validSequence(self, longer: str, shorter: str) -> List[int]:
        def calc(s: str, t: str) -> List[int]:
            res = [0] * (len(s) + 1)
            i = j = 0
            while i < len(s) and j < len(t):
                if s[i] == t[j]:
                    j += 1
                i += 1
                res[i] = j
            res[i + 1 :] = [j] * (len(s) - i)
            return res

        preC = calc(longer, shorter)
        sufC = calc(longer[::-1], shorter[::-1])[::-1]

        def getResIndex(i: int) -> List[int]:
            nonlocal shorter
            pre, suf = preC[i], sufC[i + 1]
            ptr, j = 0, 0
            res = []
            while len(res) < pre:
                if longer[ptr] == shorter[j]:
                    res.append(ptr)
                    j += 1
                ptr += 1
            f = longer[ptr] == shorter[j]
            res.append(ptr)
            ptr += 1
            j += 1
            while len(res) < len(shorter):
                if longer[ptr] == shorter[j] or f:
                    f &= longer[ptr] == shorter[j]
                    res.append(ptr)
                    j += 1
                ptr += 1
            return res

        for i in range(len(longer)):
            if preC[i] + sufC[i + 1] + 1 >= len(shorter):
                return getResIndex(i)
        return []

        print(calc(longer, shorter))
        print(calc(longer[::-1], shorter[::-1]))


# word1 = "vbcca", word2 = "abc"
# word1 = "bacdc", word2 = "abc"
# print(Solution().validSequence("vbcca", "abc"))
# print(Solution().validSequence("bacdc", "abc"))
# # # word1 = "aaaaaa", word2 = "aaabc"
# # print(Solution().validSequence("aaaaaa", "aaabc"))
# # word1 = "abc", word2 = "ab"
# print(Solution().validSequence("abc", "ab"))
# "cbbccc"
# "bb"
# print(Solution().validSequence("cbbccc", "bb"))
# "cbcbccbbcbcbbbb"
# "bb"
# print(Solution().validSequence("cbcbccbbcbcbbbb", "bb"))
# "ghhgghhhhhh"
# "gg"
# print(Solution().validSequence("ghhgghhhhhh", "gg"))
# "cccbcccbccbcccccccb"
# "ccb"
print(Solution().validSequence("cccbcccbccbcccccccb", "ccb"))
