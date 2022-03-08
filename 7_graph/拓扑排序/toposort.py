from typing import List, Tuple
from collections import defaultdict, deque


def toposort(n: int, preDeps: List[List[int]]) -> Tuple[int, List[int]]:
    """"返回拓扑排序方案数和拓扑排序结果"""
    adjMap = defaultdict(set)
    indegree = defaultdict(int)
    vertex = set(range(n))
    visitedPair = set()
    for pre, next in preDeps:
        # 保证indegree不重复计算
        if (pre, next) not in visitedPair:
            visitedPair.add((pre, next))
            adjMap[pre].add(next)
            indegree[next] += 1
            # vertex |= {pre, next}

    queue = deque([v for v in vertex if indegree[v] == 0])
    res = []
    topoCount = 1
    while queue:
        topoCount *= len(queue)
        cur = queue.popleft()
        res.append(cur)
        for next in adjMap[cur]:
            indegree[next] -= 1
            if indegree[next] == 0:
                queue.append(next)

    if len(res) != n:
        return 0, []

    return topoCount, res


print(toposort(5, [[1, 4], [2, 4], [3, 1], [3, 2]]))
