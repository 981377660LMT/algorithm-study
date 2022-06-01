# 广义邻居
from collections import defaultdict, deque
from typing import List


def genNexts(word: str):
    yield from (word[:i] + '*' + word[i + 1 :] for i in range(len(word)))


class Solution:
    def minMutation(self, start: str, end: str, bank: List[str]) -> int:
        adjMap = defaultdict(set)
        for cur in bank:
            for mid in genNexts(cur):
                adjMap[mid].add(cur)

        queue = deque([(start, 0)])
        visited = set([start])
        while queue:
            cur, dist = queue.popleft()
            if cur == end:
                return dist
            for mid in genNexts(cur):
                for next in adjMap[mid]:
                    if next in visited:
                        continue
                    visited.add(next)
                    queue.append((next, dist + 1))
        return -1


print(
    Solution().minMutation("AACCTTGG", "AATTCCGG", ["AATTCCGG", "AACCTGGG", "AACCCCGG", "AACCTACC"])
)

