from typing import List
from collections import deque

# TSP problem
# 1 <= n <= 12
class Solution:
    def shortestPathLength(self, graph: List[List[int]]) -> int:
        n = len(graph)
        target = (1 << n) - 1
        visited = [[0x7FFFFFFF] * (1 << n) for _ in range(n)]
        queue = deque([(i, 1 << i, 0, str(i)) for i in range(n)])

        while queue:
            cur, state, cost, path = queue.popleft()
            if state == target:
                print(list(path))
                return cost
            for next in graph[cur]:
                nextState = state | (1 << next)
                if visited[next][nextState] > cost + 1:
                    visited[next][nextState] = cost + 1
                    queue.append((next, nextState, cost + 1, path + str(next)))

        return -1


print(Solution().shortestPathLength(graph=[[1, 2, 3], [0], [0], [0]]))
# 输出：4
# 解释：一种可能的路径为 [1,0,2,0,3]
