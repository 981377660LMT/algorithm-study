from typing import List, Tuple


class BipartiteMathing:
    """二分图最大匹配."""

    __slots__ = "_timestamp", "_graph", "_alive", "_used", "_match", "_n"

    def __init__(self, n: int) -> None:
        self._timestamp = 0
        self._graph = [[] for _ in range(n)]
        self._alive = [True] * n
        self._used = [0] * n
        self._match = [-1] * n
        self._n = n

    def addEdge(self, left: int, right: int) -> None:
        """添加一条边.

        Args:
            left (int): 左侧顶点.0 <= left < x.
            right (int): 右侧顶点.x <= right < x + y (n).
        """
        self._graph[left].append(right)
        self._graph[right].append(left)

    def removeEdge(self, u: int, v: int) -> None:
        self._graph[u].remove(v)
        self._graph[v].remove(u)

    def maxMatching(self) -> List[Tuple[int, int]]:
        """!O(VE)求二分图最大匹配."""
        n = self._n
        self._match = [-1] * n  # 重置匹配
        alive, match = self._alive, self._match
        for u in range(n):
            if alive[u] and match[u] == -1:
                self._timestamp += 1
                self._argument(u)
        return [(u, match[u]) for u in range(n) if u < match[u]]

    def removeVertex(self, idx: int) -> int:
        """删除顶点idx, 返回流量的变化量(-1/0)."""
        alive, match = self._alive, self._match
        alive[idx] = False
        if match[idx] == -1:
            return 0
        match[match[idx]] = -1
        self._timestamp += 1
        res = self._argument(match[idx])
        match[idx] = -1
        return 0 if res else -1

    def addVertex(self, idx: int) -> int:
        """添加顶点idx, 返回流量的变化量(0/1)."""
        self._alive[idx] = True
        self._timestamp += 1
        return 1 if self._argument(idx) else 0

    def getMatchingEdges(self) -> List[Tuple[int, int]]:
        """获取匹配边.需要先调用 maxMatching."""
        return [(u, self._match[u]) for u in range(self._n) if u < self._match[u]]

    def _argument(self, idx: int) -> bool:
        used, match, graph, alive = self._used, self._match, self._graph, self._alive
        used[idx] = self._timestamp
        for to in graph[idx]:
            toMatch = match[to]
            if not alive[to]:
                continue
            if toMatch == -1 or (used[toMatch] != self._timestamp and self._argument(toMatch)):
                match[idx] = to
                match[to] = idx
                return True
        return False


def createBipartiteMathingFromEdges(
    n: int, edges: List[Tuple[int, int]]
) -> Tuple["BipartiteMathing", List[int], List[int]]:
    """从边创建二分图最大匹配.

    Args:
        n (int): 顶点数.
        edges (List[Tuple[int, int]]): 边集.

    Returns:
        Tuple[BipartiteMathing, List[int], List[int]]:
        二分图最大匹配.
        原图中的点在二分图中的编号.左侧点编号为0-L-1, 右侧点编号为L-n-1.
        二分图中的点在原图中的编号.0-n-1 -> 0-n-1.

    ```python
    G, ids, rids = createBipartiteMathingFromEdges(4, [(0, 1), (1, 2), (2, 3)])
    M = G.maxMatching()
    edges = [(rids[u], rids[v]) for u, v in M]
    print(edges) # [(0, 1), (2, 3)]
    ```
    """
    graph = [[] for _ in range(n)]
    for u, v in edges:
        graph[u].append(v)
        graph[v].append(u)
    colors, ok = isBipartite(n, graph)
    if not ok:
        raise Exception("not bipartite")

    ids, rids = [0] * n, [0] * n
    L = n - sum(colors)
    left, right = 0, 0  # 规定左侧点颜色为0, 右侧点颜色为1
    for i in range(n):
        if colors[i] == 0:
            ids[i] = left
            rids[left] = i
            left += 1
        else:
            ids[i] = right + L
            rids[right + L] = i
            right += 1

    bm = BipartiteMathing(n)
    for u, v in edges:
        if colors[u] == 1:
            u, v = v, u
        bm.addEdge(ids[u], ids[v])
    return bm, ids, rids


def isBipartite(n: int, graph: List[List[int]]) -> Tuple[List[int], bool]:
    """判断是否是二分图.返回 (染色的01数组, 是否是二分图)."""
    color = [-1] * n
    for i in range(n):
        if color[i] == -1:
            color[i] = 0
            stack = [i]
            while stack:
                cur = stack.pop()
                for next_ in graph[cur]:
                    if color[next_] == -1:
                        color[next_] = 1 ^ color[cur]
                        stack.append(next_)
                    elif color[next_] == color[cur]:
                        return [], False
    return color, True


if __name__ == "__main__":
    # https://atcoder.jp/contests/abc317/tasks/abc317_g
    # 跑m次匈牙利，每跑一次就删去完美匹配的边
    def rearrange(grid: List[List[int]]) -> List[List[int]]:
        ROW, COL = len(grid), len(grid[0])
        res = [[0] * COL for _ in range(ROW)]
        G = BipartiteMathing(ROW + ROW)
        for i, row in enumerate(grid):
            for v in row:
                G.addEdge(i, ROW + v)

        for c in range(COL):
            matching = G.maxMatching()
            for u, v in matching:
                res[u][c] = v - ROW + 1
                G.removeEdge(u, v)
        return res

    import sys

    input = sys.stdin.readline
    row, col = map(int, input().split())
    grid = []
    for _ in range(row):
        row = [x - 1 for x in map(int, input().split())]
        grid.append(row)
    res = rearrange(grid)

    print("Yes")
    for row in res:
        print(*row)
