from typing import List, Tuple
from collections import defaultdict, Counter, deque

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def hasCycle(self, graph: str) -> bool:
        words = graph.split(',')
        adjMap = defaultdict(set)
        deg = defaultdict(int)
        allVertex = set()
        for w in words:
            a, b = w.split('->')
            adjMap[a].add(b)
            deg[b] += 1
            allVertex.add(a)
            allVertex.add(b)
        queue = deque([i for i in allVertex if deg[i] == 0])
        while queue:
            cur = queue.popleft()
            for n in adjMap[cur]:
                deg[n] -= 1
                if deg[n] == 0:
                    queue.append(n)
        return any(deg[i] != 0 for i in allVertex)

