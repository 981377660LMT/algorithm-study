from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个数组 arr ，数组中有 n 个 非空 字符串。

# 请你求出一个长度为 n 的字符串 answer ，满足：


# answer[i] 是 arr[i] 最短 的子字符串，且它不是 arr 中其他任何字符串的子字符串。如果有多个这样的子字符串存在，answer[i] 应该是它们中字典序最小的一个。如果不存在这样的子字符串，answer[i] 为空字符串。
# 请你返回数组 answer 。
BIG = chr(127) * 30


class Solution:
    def shortestSubstrings(self, arr: List[str]) -> List[str]:
        res = [""] * len(arr)
        for i, s in enumerate(arr):
            other = arr[:]
            other.pop(i)
            t = "#".join(other)
            cand = BIG
            for a in range(len(s)):
                for b in range(a + 1, len(s) + 1):
                    tmp = s[a:b]
                    if tmp not in t:
                        if len(tmp) < len(cand) or (len(tmp) == len(cand) and tmp < cand):
                            cand = tmp
            if cand != BIG:
                res[i] = cand
        return res


# ["gfnt","xn","mdz","yfmr","fi","wwncn","hkdy"]


print(Solution().shortestSubstrings(["gfnt", "xn", "mdz", "yfmr", "fi", "wwncn", "hkdy"]))
