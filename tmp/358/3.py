from itertools import permutations
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你三个字符串 a ，b 和 c ， 你的任务是找到长度 最短 的字符串，且这三个字符串都是它的 子字符串 。
# 如果有多个这样的字符串，请你返回 字典序最小 的一个。

# 请你返回满足题目要求的字符串。

# 注意：


# 两个长度相同的字符串 a 和 b ，如果在第一个不相同的字符处，a 的字母在字母表中比 b 的字母 靠前 ，那么字符串 a 比字符串 b 字典序小 。
# 子字符串 是一个字符串中一段连续的字符序列。


class Solution:
    def minimumString(self, a: str, b: str, c: str) -> str:
        def maxCommon(pre: str, post: str) -> int:
            """pre的后缀和post的前缀的最大公共长度"""
            res = 0
            for i in range(1, len(pre) + 1):
                if pre[-i:] == post[:i]:
                    res = i
            return res

        res = []
        for perm in permutations([a, b, c]):
            w1, w2, w3 = perm
            if w2 not in w1:
                common1 = maxCommon(w1, w2)
                w1 = w1 + w2[common1:]
            if w3 not in w1:
                common2 = maxCommon(w1, w3)
                w1 = w1 + w3[common2:]
            res.append(w1)
        res.sort(key=lambda x: (len(x), x))
        return res[0]


# print(Solution().minimumString(a="abc", b="bca", c="aaa"))
# "ca"
# "a"
# "a"
print(Solution().minimumString(a="ca", b="a", c="a"))  # ca
# "cab"
# "a"
# "b"
print(Solution().minimumString(a="cab", b="a", c="a"))  # cab
