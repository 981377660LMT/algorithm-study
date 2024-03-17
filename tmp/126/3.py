from itertools import groupby
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个字符串 s 。s[i] 要么是小写英文字母，要么是问号 '?' 。

# 对于长度为 m 且 只 含有小写英文字母的字符串 t ，我们定义函数 cost(i) 为下标 i 之前（也就是范围 [0, i - 1] 中）出现过与 t[i] 相同 字符出现的次数。

# 字符串 t 的 分数 为所有下标 i 的 cost(i) 之 和 。

# 比方说，字符串 t = "aab" ：

# cost(0) = 0
# cost(1) = 1
# cost(2) = 0
# 所以，字符串 "aab" 的分数为 0 + 1 + 0 = 1 。
# 你的任务是用小写英文字母 替换 s 中 所有 问号，使 s 的 分数最小 。


# 请你返回替换所有问号 '?' 之后且分数最小的字符串。如果有多个字符串的 分数最小 ，那么返回字典序最小的一个。
class Solution:
    def minimizeStringValue(self, s: str) -> str:
        # 每次贪心选择出现次数最少的最字符.如果出现次数相同,选择字典序最小的字符.
        # 每一段字典序最小???
        counter = [0] * 26
        res = []
        groups = [(char, len(list(group))) for char, group in groupby(s)]
        for char, count in groups:
            if char != "?":
                counter[ord(char) - 97] += count

        ptr = 0
        stash = []
        for char, count in groups[::-1]:
            if char != "?":
                # counter[ord(char) - 97] += count
                ...
            else:
                curGroup = []
                for _ in range(ptr, ptr + count):
                    min_ = min(counter)
                    minIndex = counter.index(min_)
                    curGroup.append(chr(minIndex + 97))
                    counter[minIndex] += 1
                stash.extend(curGroup)
            ptr += count

        res = []
        stash.sort(reverse=True)
        for char, count in groups:
            if char != "?":
                res.append(char * count)
            else:
                for _ in range(count):
                    res.append(stash.pop())

        return "".join(res)


# "abcdefghijklmnopqrstuvwxy??"
print(Solution().minimizeStringValue("abcdefghijklmnopqrstuvwxy??"))
# "eq?umjlasi"
print(Solution().minimizeStringValue("eq?umjlasi"))
# "g?xvgroui??xk?zqb?da?jan?cdhtksme"
print(Solution().minimizeStringValue("g?xvgroui??xk?zqb?da?jan?cdhtksme"))
