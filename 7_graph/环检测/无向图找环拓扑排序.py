from collections import defaultdict, deque
from typing import DefaultDict, List, Set, Tuple


def findCycle(n: int, edges: List[Tuple[int, int]]) -> List[int]:
    """无向图找环上的点"""
    adjMap = defaultdict(set)
    degrees = [0] * n
    for u, v in edges:
        adjMap[u].add(v)
        adjMap[v].add(u)
        degrees[v] += 1
        degrees[u] += 1

    queue = deque([i for i in range(n) if degrees[i] == 1])
    onCycle = [True] * n
    while queue:
        cur = queue.popleft()
        onCycle[cur] = False
        for next in adjMap[cur]:
            degrees[next] -= 1
            if degrees[next] == 1:
                queue.append(next)

    cycle = [i for i, v in enumerate(onCycle) if v]
    return cycle

