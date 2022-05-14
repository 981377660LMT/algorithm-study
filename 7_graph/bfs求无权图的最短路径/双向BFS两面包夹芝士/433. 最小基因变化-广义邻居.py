# 广义邻居
from collections import defaultdict, deque
from typing import List


class Solution:
    def minMutation(self, start: str, end: str, bank: List[str]) -> int:
        adjMap = defaultdict(set)
        for cur in bank:
            for i in range(len(cur)):
                next = cur[:i] + '*' + cur[i + 1 :]
                adjMap[next].add(cur)

        queue = deque([(start, 0)])
        visited = set([start])
        while queue:
            cur, dist = queue.popleft()
            if cur == end:
                return dist
            for i in range(len(cur)):
                for next in adjMap[cur[:i] + '*' + cur[i + 1 :]]:
                    if next in visited:
                        continue
                    visited.add(next)
                    queue.append((next, dist + 1))
        return -1


print(
    Solution().minMutation("AACCTTGG", "AATTCCGG", ["AATTCCGG", "AACCTGGG", "AACCCCGG", "AACCTACC"])
)

