# You are given a string s containing 1s and 0s.
# Given that you can swap at most one pair of characters in the string,
# return the resulting length of the longest contiguous substring of 1s.
from collections import Counter
from itertools import groupby


class Solution:
    def solve(self, s: str):
        counter = Counter(s)
        groups = [[char, len(list(group))] for char, group in groupby(s)]

        # 1. extend 1 情形
        res = max(
            (count + int(counter[char] > count) for char, count in groups if char == '1'), default=0
        )

        # 2. remove 1 divider 情形
        for i in range(1, len(groups) - 1):
            if groups[i - 1][0] == groups[i + 1][0] == '1' and groups[i][1] == 1:
                curSum = groups[i - 1][1] + groups[i + 1][1] + 1
                res = max(res, curSum)

        return min(res, counter['1'])


print(Solution().solve(s="10101"))
print(Solution().solve(s="11011"))
