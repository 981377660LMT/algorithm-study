from typing import List, Tuple
from collections import defaultdict, deque, Counter
from heapq import heapify, heappop, heappush
from sortedcontainers import SortedList, SortedDict
from bisect import bisect_left, bisect_right
from functools import lru_cache, reduce
from itertools import accumulate, groupby, combinations, permutations, product
from math import gcd, sqrt, ceil, floor, comb

MOD = int(1e9 + 7)
INF = 0x3F3F3F3F
EPS = int(1e-6)
dirs4 = [[-1, 0], [0, 1], [1, 0], [0, -1]]
dirs8 = [[-1, 0], [-1, 1], [0, 1], [1, 1], [1, 0], [1, -1], [0, -1], [-1, -1]]


# class Solution:
#     def longestPalindrome(self, words: List[str]) -> int:
#         counter = Counter(words)
#         print(counter)
#         res = 0
#         same = Counter()
#         visited = set()
#         ok = set()

#         for w in words:
#             if w[0] == w[1]:
#                 same[w] += 1
#             elif w[::-1] in visited:
#                 res += 4
#                 visited.discard(w[::-1])
#                 ok.add(w[0])
#                 ok.add(w[1])
#             else:
#                 visited.add(w)
#         # print(same, res)

#         for w, count in list(same.items()):
#             pair = count // 2
#             res += pair * 4
#             same[w] -= pair * 2
#             if not same[w]:
#                 del same[w]
#         print(same, res)

#         # hit = 0
#         # for w, count in same.items():
#         #     if w[0] in ok:
#         #         res += count * 2
#         #         hit += 1

#         # if hit != len(same):
#         res += min(2, 2 * len(same))

#         return res


# counter配对题
# 这题卡了很久，是因为还没想清楚就开始写了，然后debug好久；以后一定要先把情况考虑清楚，丢到counter里辅助思考

# 情况就是配对+一个same
class Solution:
    def longestPalindrome(self, words: List[str]) -> int:
        same = Counter()
        diff = Counter()

        for word in words:
            if word[0] == word[1]:
                same[word] += 1
            else:
                diff[word] += 1

        res = 0
        for word in diff:
            match = word[::-1]
            if match in diff:
                res += min(diff[word], diff[match]) * 2
        for word in same:
            res += same[word] // 2 * 4
            same[word] %= 2

        for word in same:
            if same[word] == 1:
                res += 2
                break

        return res


# 6 8 2 14 22 14
print(Solution().longestPalindrome(words=["lc", "cl", "gg"]))
print(Solution().longestPalindrome(words=["ab", "ty", "yt", "lc", "cl", "ab"]))
print(Solution().longestPalindrome(words=["cc", "ll", "xx"]))
print(
    Solution().longestPalindrome(
        words=[
            "qo",
            "fo",
            "fq",
            "qf",
            "fo",
            "ff",
            "qq",
            "qf",
            "of",
            "of",
            "oo",
            "of",
            "of",
            "qf",
            "qf",
            "of",
        ]
    )
)
print(
    Solution().longestPalindrome(
        ["dd", "aa", "bb", "dd", "aa", "dd", "bb", "dd", "aa", "cc", "bb", "cc", "dd", "cc"]
    )
)
print(
    Solution().longestPalindrome(words=["em", "pe", "mp", "ee", "pp", "me", "ep", "em", "em", "me"])
)
