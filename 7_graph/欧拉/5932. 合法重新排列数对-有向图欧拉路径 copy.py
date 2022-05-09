from typing import List
from collections import defaultdict
from getEulerPath import getEulerPath


class Solution:
    def validArrangement(self, pairs: List[List[int]]) -> List[List[int]]:
        adjMap = defaultdict(set)
        allVertex = set()
        for u, v in pairs:
            adjMap[u].add(v)
            allVertex |= {u, v}
        path = getEulerPath(allVertex, adjMap, isDirected=True)[1]
        return [[pre, cur] for pre, cur in zip(path, path[1:])]


print(Solution().validArrangement(pairs=[[5, 1], [4, 5], [11, 9], [9, 4]]))

