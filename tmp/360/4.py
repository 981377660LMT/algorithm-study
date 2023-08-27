from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个长度为 n 下标从 0 开始的整数数组 receiver 和一个整数 k 。

# 总共有 n 名玩家，玩家 编号 互不相同，且为 [0, n - 1] 中的整数。这些玩家玩一个传球游戏，receiver[i] 表示编号为 i 的玩家会传球给编号为 receiver[i] 的玩家。玩家可以传球给自己，也就是说 receiver[i] 可能等于 i 。

# 你需要从 n 名玩家中选择一名玩家作为游戏开始时唯一手中有球的玩家，球会被传 恰好 k 次。

# 如果选择编号为 x 的玩家作为开始玩家，定义函数 f(x) 表示从编号为 x 的玩家开始，k 次传球内所有接触过球玩家的编号之 和 ，如果有玩家多次触球，则 累加多次 。换句话说， f(x) = x + receiver[x] + receiver[receiver[x]] + ... + receiver(k)[x] 。

# 你的任务时选择开始玩家 x ，目的是 最大化 f(x) 。

# 请你返回函数的 最大值 。


# 注意：receiver 可能含有重复元素。


# 枚举起始点

from collections import deque
from typing import List, Tuple


def findCycleAndCalDepth(
    n: int, adjList: List[List[int]], deg: List[int], *, isDirected: bool
) -> Tuple[List[List[int]], List[int]]:
    """无/有向基环树森林找环上的点,并记录每个点在拓扑排序中的最大深度,最外层的点深度为0"""

    def max(a: int, b: int) -> int:
        return a if a > b else b

    depth = [0] * n
    startDeg = 0 if isDirected else 1
    queue = deque([i for i in range(n) if deg[i] == startDeg])
    visited = [False] * n
    while queue:
        cur = queue.popleft()
        visited[cur] = True
        for next in adjList[cur]:
            depth[next] = max(depth[next], depth[cur] + 1)
            deg[next] -= 1
            if deg[next] == startDeg:
                queue.append(next)

    return cycleGroup, depth


class Solution:
    def getMaxFunctionValue(self, receiver: List[int], k: int) -> int:
        n = len(receiver)
        adjList = [[] for _ in range(n)]
        for i in range(n):
            next_ = receiver[i]
            adjList[i].append(next_)

        visited = [False] * n

        def dfs(cur: int, path: List[int]) -> None:
            if visited[cur]:
                return
            visited[cur] = True
            path.append(cur)
            for next in adjList[cur]:
                dfs(next, path)

        cycleGroup = []
        for i in range(n):
            if visited[i]:
                continue
            path = []
            dfs(i, path)
            cycleGroup.append(path)

        print(cycleGroup)


# receiver = [2,0,1], k = 4
print(Solution().getMaxFunctionValue([2, 0, 1], 4))
# receiver = [1,1,1,2,3], k = 3
print(Solution().getMaxFunctionValue([1, 1, 1, 2, 3], 3))
