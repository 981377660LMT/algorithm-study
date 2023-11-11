# https://atcoder.jp/contests/abc241/tasks/abc241_g

# 巡回对战(総当たり戦/round robin )的橄榄球比赛
# n个队每个队要打n-1场比赛 每场比赛必须分出胜负
# 现在已知m场比赛的结果
# !对每个人x判断是否可能成为胜者(胜利场数最多且唯一)。
# n<=50

# https://www.bilibili.com/read/cv15437330?spm_id_from=333.1007.0.0

# 橄榄球比赛建模(Baseball elimination) => 最大流算法
# !对每个人，看他是否能够取胜
# !建图，有O(n^2)条边 容量为O(n) 所以Dinic每次查询复杂度不超过O(maxflow*E) 即 O(n^3)

# 源点 => 分配胜利给每个场次
# 每个场次 => 将胜利分给两队中的一个
# 每个队最多还可以赢几场比赛 => 流向汇点
# !如果最大流等于n*(n-1) 那么就可能成为获胜者


import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")


from collections import defaultdict, deque
from typing import Set

INF = int(4e18)


class ATCMaxFlow:
    """Dinic算法 数组+边存图"""

    __slots__ = (
        "_n",
        "_start",
        "_end",
        "_reGraph",
        "_edges",
        "_visitedEdge",
        "_levels",
        "_curEdges",
    )

    def __init__(self, n: int, start: int, end: int) -> None:
        if not (0 <= start < n and 0 <= end < n):
            raise ValueError(f"start: {start}, end: {end} out of range [0,{n}]")

        self._n = n
        self._start = start
        self._end = end
        self._reGraph = [[] for _ in range(n)]  # 残量图存边的序号
        self._edges = []  # [next,capacity]

        self._visitedEdge = set()

        self._levels = [0] * n
        self._curEdges = [0] * n

    def addEdge(self, v1: int, v2: int, capacity: int) -> None:
        """添加边 v1->v2, 容量为w 注意会添加重边"""
        self._visitedEdge.add((v1, v2))
        self._reGraph[v1].append(len(self._edges))
        self._edges.append([v2, capacity])
        self._reGraph[v2].append(len(self._edges))
        self._edges.append([v1, 0])

    def addEdgeIfAbsent(self, v1: int, v2: int, capacity: int) -> None:
        """如果边不存在则添加边 v1->v2, 容量为w"""
        if (v1, v2) in self._visitedEdge:
            return
        self._visitedEdge.add((v1, v2))
        self._reGraph[v1].append(len(self._edges))
        self._edges.append([v2, capacity])
        self._reGraph[v2].append(len(self._edges))
        self._edges.append([v1, 0])

    def calMaxFlow(self) -> int:
        n, start, end = self._n, self._start, self._end
        res = 0

        while self._bfs():
            self._curEdges = [0] * n
            res += self._dfs(start, end, INF)
        return res

    def getPath(self) -> Set[int]:
        """最大流经过了哪些点"""
        visited = set()
        queue = [self._start]
        reGraph, edges = self._reGraph, self._edges
        while queue:
            cur = queue.pop()
            visited.add(cur)
            for ei in reGraph[cur]:
                edge = edges[ei]
                next, remain = edge
                if remain > 0 and next not in visited:
                    visited.add(next)
                    queue.append(next)
        return visited

    def useQueryRemainOfEdge(self):
        """求边的残量(剩余的容量)::

        ```python
        maxFlow = ATCMaxFlow(n, start, end)
        query = maxFlow.useQueryRemainOfEdge()
        edgeRemain = query(v1, v2)
        ```
        """

        def query(v1: int, v2: int) -> int:
            return adjList[v1][v2]

        n, reGraph, edges = self._n, self._reGraph, self._edges
        adjList = [defaultdict(int) for _ in range(n)]
        for cur in range(n):
            for ei in reGraph[cur]:
                edge = edges[ei]
                next, remain = edge
                adjList[cur][next] += remain

        return query

    def _bfs(self) -> bool:
        n, reGraph, start, end, edges = self._n, self._reGraph, self._start, self._end, self._edges
        self._levels = level = [-1] * n
        level[start] = 0
        queue = deque([start])

        while queue:
            cur = queue.popleft()
            nextDist = level[cur] + 1
            for ei in reGraph[cur]:
                next, remain = edges[ei]
                if remain > 0 and level[next] == -1:
                    level[next] = nextDist
                    if next == end:
                        return True
                    queue.append(next)

        return False

    def _dfs(self, cur: int, end: int, flow: int) -> int:
        if cur == end:
            return flow
        res = flow
        reGraph, level, curEdges, edges = self._reGraph, self._levels, self._curEdges, self._edges
        ei = curEdges[cur]
        while ei < len(reGraph[cur]):
            ej = reGraph[cur][ei]
            next, remain = edges[ej]
            if remain > 0 and level[cur] + 1 == level[next]:
                delta = self._dfs(next, end, min(res, remain))
                edges[ej][1] -= delta
                edges[ej ^ 1][1] += delta
                res -= delta
                if res == 0:
                    return flow
            curEdges[cur] += 1
            ei = curEdges[cur]

        return flow - res


def query(person: int) -> bool:
    """person是否有胜利的可能

    假设他之后全胜 别人必须比他的得分严格小
    如果最大流等于总场数 那么存在分配方案
    即person可能成为胜者
    """
    v = 2500
    START, END, OFFSET = 2 * v, 2 * v + 1, v
    maxFlow = ATCMaxFlow(2 * v + 2, START, END)

    # match[i][j]表示i和j之间的比赛
    # 如果i获胜 连接(game,i+OFFSET)
    # 如果j获胜 连接(game,j+OFFSET)
    # 如果没有比过 连接(game,i+OFFSET)和(game,j+OFFSET)
    for i in range(n):
        for j in range(i + 1, n):
            game = i * n + j
            maxFlow.addEdgeIfAbsent(START, game, 1)
            if match[i][j] == -1:
                maxFlow.addEdgeIfAbsent(game, i + OFFSET, 1)
                maxFlow.addEdgeIfAbsent(game, j + OFFSET, 1)
            elif match[i][j] == i:
                maxFlow.addEdgeIfAbsent(game, i + OFFSET, 1)
            else:
                maxFlow.addEdgeIfAbsent(game, j + OFFSET, 1)

    maxWin = win[person] + remain[person]
    for i in range(n):
        if i == person:
            maxFlow.addEdgeIfAbsent(i + OFFSET, END, maxWin)
        else:
            maxFlow.addEdgeIfAbsent(i + OFFSET, END, maxWin - 1)  # 别的队最多赢几场比赛

    res = maxFlow.calMaxFlow()
    return res == n * (n - 1) // 2  # 最大流等于总场数


n, m = map(int, input().split())
match = [[-1] * n for _ in range(n)]  # !对战记录
win, lose, remain = [0] * n, [0] * n, [n - 1] * n

for _ in range(m):
    u, v = map(int, input().split())  # !u赢了v
    u, v = u - 1, v - 1
    match[u][v] = match[v][u] = u
    win[u] += 1
    lose[v] += 1
    remain[u] -= 1
    remain[v] -= 1


res = []
for i in range(n):
    if query(i):
        res.append(i + 1)
print(*res)
