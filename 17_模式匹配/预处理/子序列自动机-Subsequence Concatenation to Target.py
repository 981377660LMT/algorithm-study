# Given strings source and target,
# return the minimum number of subsequences of source we can form such that their concatenation equals target.
# If it's not possible, return -1.
# n ≤ 100,000
from typing import List, Tuple

# 连接子串获得目标串的最小需要量
# n ≤ 100,000


class Solution:
    def solve(self, source: str, target: str):
        n = len(source)
        nexts: List[Tuple[int, ...]] = [()] * (n + 1)
        last = [-1] * 26
        nexts[-1] = tuple(last)
        for i in range(n - 1, -1, -1):
            last[ord(source[i]) - 97] = i  # 注意这里先更新
            nexts[i] = tuple(last)

        # print(nexts)

        searchFrom = 0
        res = 0
        for char in target:
            ord_ = ord(char) - 97
            if nexts[0][ord_] == -1:
                return -1

            nextIndex = nexts[searchFrom][ord_]
            if nextIndex == -1:
                res += 1
                nextIndex = nexts[0][ord_]
            searchFrom = nextIndex + 1

        return res + 1


print(Solution().solve(source="abc", target="abcbcc"))
# We can have the following subsequences of source: "abc" + "bc" + "c"
