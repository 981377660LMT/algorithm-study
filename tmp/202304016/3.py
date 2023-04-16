from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个字符串 word ，你可以向其中任何位置插入 "a"、"b" 或 "c" 任意次，返回使 word 有效 需要插入的最少字母数。

# TODO 完全不会的贪心


# 如果字符串可以由 "abc" 串联多次得到，则认为该字符串 有效 。
class Solution:
    def addMinimum(self, word: str) -> int:
        # remove abc
        pre = ""
        res = 0
        for c in word:
            if not pre:
                pre = c
                continue

            if c == "a":
                if pre == "a":
                    res += 2
                elif pre == "b":
                    res += 2
                elif pre == "c":
                    res += 2
                elif pre == "ab":
                    res += 1
                elif pre == "ac":
                    res += 1
                elif pre == "bc":
                    res += 1
                pre = c
            elif c == "b":
                if pre == "a":
                    pre = "ab"
                    continue
                elif pre == "b":
                    res += 2
                elif pre == "c":
                    res += 2
                elif pre == "ab":
                    res += 1
                elif pre == "ac":
                    res += 1
                elif pre == "bc":
                    res += 1
                pre = c
            else:
                if pre == "a":
                    pre = "ac"
                    continue
                elif pre == "b":
                    pre = "bc"
                    continue
                elif pre == "c":
                    res += 2
                elif pre == "ab":
                    pre = ""
                    continue
                elif pre == "ac":
                    res += 1
                elif pre == "bc":
                    res += 1
                pre = c

        if pre:
            res += 3 - len(pre)
        return res


print(Solution().addMinimum("abcabc"))
# "aaaaac"
print(Solution().addMinimum("aaaaac"))
