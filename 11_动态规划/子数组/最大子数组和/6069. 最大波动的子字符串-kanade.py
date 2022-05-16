"""求 子串中出现次数 最多 的字符次数与出现次数 最少 的字符次数之差 的 最大值"""
import gc
from itertools import combinations


gc.disable()


class Solution:
    def largestVariance(self, s: str) -> int:
        """时间复杂度O(26*26*n)"""

        def cal(s1: str, s2: str) -> int:
            # s1 最多 s2 最少 注意必须包含 s2
            # 用一个变量记录包含s2时的差值
            res = 0
            maxSum, maxSumWithS2 = 0, -int(1e9)  # 一开始没有s2
            for char in s:
                if char == s1:
                    maxSum += 1
                    maxSumWithS2 += 1
                elif char == s2:
                    maxSum -= 1
                    maxSumWithS2 = maxSum  # 注意这里更新最大值
                    if maxSum < 0:  # 全不要了
                        maxSum = 0
                if maxSumWithS2 > res:
                    res = maxSumWithS2  # if 代替 max 2000+ms
                # res = max(res, maxSumWithS2)  # 用 max 7000+ms
            return res

        allChar = list(set(s))
        res = 0
        for s1, s2 in combinations(allChar, 2):
            res = max(res, cal(s1, s2), cal(s2, s1))
        return res

    def largestVariance2(self, s: str) -> int:
        """时间复杂度O(26*n)"""

        res = 0
        sumMax = [[0] * 26 for _ in range(26)]
        sumMaxWiths2 = [[-int(1e20)] * 26 for _ in range(26)]
        for char in s:
            ord_ = ord(char) - ord('a')
            for i in range(26):
                if i == ord_:
                    continue
                sumMax[ord_][i] += 1
                sumMaxWiths2[ord_][i] += 1
                sumMax[i][ord_] -= 1
                sumMaxWiths2[i][ord_] = sumMax[i][ord_]
                if sumMax[i][ord_] < 0:
                    sumMax[i][ord_] = 0
                res = max(res, sumMaxWiths2[ord_][i], sumMaxWiths2[i][ord_])
        return res


print(Solution().largestVariance("aababbb"))
print(Solution().largestVariance("abcde"))
print(Solution().largestVariance("bbaabbaabbaabb"))  # 2
print(Solution().largestVariance("lripaa"))  # 2
# from itertools import combinations
# from string import ascii_lowercase
# from typing import List, Tuple
# from collections import defaultdict

# MOD = int(1e9 + 7)
# INF = int(1e20)


# class AlphaPresum:
#     def __init__(self, s: str) -> None:
#         preSum = [[0] * 26]
#         for char in s:
#             cur = preSum[-1][:]
#             cur[ord(char) - ord('a')] += 1
#             preSum.append(cur)

#         self._preSum = preSum

#     def getCountOfSlice(self, char: str, left: int, right: int) -> int:
#         """
#         >>> preSum = AlphaPresum("abcdabcd")
#         >>> print(preSum.getCountOfSlice('a', 0, 2)) # s[0:2]间'a'个数为1
#         >>> print(preSum.getCountOfSlice('a', 0, 8)) # s[0:8]间'a'个数为2
#         """
#         assert 0 <= left <= right < len(self._preSum)
#         ord_ = ord(char) - ord('a')
#         return self._preSum[right][ord_] - self._preSum[left][ord_]


# class Solution:
#     def largestVariance(self, s: str) -> int:
#         n = len(s)
#         allChar = list(set(s))
#         preSum = AlphaPresum(s)

#         adjMap = defaultdict(list)

#         for i, char in enumerate(s):
#             adjMap[char].append(i)

#         # 每个数作为最小值
#         res = 0
#         for c in allChar:
#             indexes = [0] + adjMap[c] + [n]
#             for i in range(1, len(indexes) - 1):
#                 pre, cur = indexes[i - 1], indexes[i + 1]
#                 for other in allChar:
#                     if c == other:
#                         continue
#                     count = preSum.getCountOfSlice(other, pre, cur)
#                     res = max(res, count - 1)

#         return res


# print(Solution().largestVariance("aababbb"))
# print(Solution().largestVariance("abcde"))
# print(Solution().largestVariance("bbaabbaabbaabb"))
