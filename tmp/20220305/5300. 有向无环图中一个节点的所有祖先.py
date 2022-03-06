from collections import defaultdict, deque
from typing import List, Set, Tuple

MOD = int(1e9 + 7)

# 有向无环图DAG必有拓扑序


class Solution:
    def getAncestors(self, n: int, edges: List[List[int]]) -> List[List[int]]:
        adjMap = defaultdict(set)
        indegree = [0] * n
        for edge in edges:
            adjMap[edge[0]].add(edge[1])
            indegree[edge[1]] += 1

        queue = deque([i for i in range(n) if indegree[i] == 0])
        res: List[Set[int]] = [set() for _ in range(n)]
        while queue:
            cur = queue.popleft()
            for next in adjMap[cur]:
                indegree[next] -= 1
                res[next] |= res[cur] | {cur}
                if indegree[next] == 0:
                    queue.append(next)
        return [sorted(s) for s in res]


[[], [], [], [0, 1], [0, 2], [0, 1, 3], [0, 1, 2, 3, 4], [0, 1, 2, 3]]
[[], [0], [0, 1], [0, 1, 2], [0, 1, 2, 3]]
