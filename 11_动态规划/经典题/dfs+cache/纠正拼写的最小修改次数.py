# n * m ≤ 200,000 where n is the length of itinerary, m is the length of edges
from collections import defaultdict
from functools import lru_cache
from typing import List

# 纠正拼写的最小修改次数，使得航线itinerary与edges逻辑正确
# 每个航班的名称为3个字母


class Solution:
    def solve(self, itinerary: List[str], edges: List[List[str]]) -> int:
        """状态：(当前航班index,之前的航班)"""

        def strDiff(s1, s2) -> int:
            return sum(1 for i in range(len(s1)) if s1[i] != s2[i])

        @lru_cache(None)
        def dfs(index: int, curPos: str) -> int:
            if index == len(itinerary):
                return 0

            target = itinerary[index]
            res = int(1e20)

            if curPos == '':
                for next in allPort:
                    res = min(res, dfs(index + 1, next) + strDiff(target, next))
            else:
                for next in adjMap[curPos]:
                    res = min(res, dfs(index + 1, next) + strDiff(target, next))

            return res

        adjMap = defaultdict(set)
        allPort = set()
        for u, v in edges:
            adjMap[u].add(v)
            allPort.add(u)

        return dfs(0, "")


# itinerary = ["YYZ", "SFO", "JFK"]
# edges = [
#     ["YYZ", "SEA"],
#     ["SEA", "JAM"],
#     ["SEA", "JFL"]
# ]
# We change "SFO" to "SEA" for 2 character changes and "JFK" to 'JFL" for 1 character change.
# In total, we made 3 character changes.
print(
    Solution().solve(
        itinerary=["YYZ", "SFO", "JFK"], edges=[["YYZ", "SEA"], ["SEA", "JAM"], ["SEA", "JFL"]],
    )
)

