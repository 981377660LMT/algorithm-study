from collections import defaultdict
from typing import Iterable, List, Mapping, Optional, Sequence, Tuple, Union


class Hungarian:
    def __init__(self, row: int, col: int):
        """匈牙利算法求无权二分图最大匹配

        时间复杂度O(V * E)

        Args:
            row (int): 男孩的个数
            col (int): 女孩的个数
        """
        self._row = row
        self._col = col
        self._graph = [[] for _ in range(row)]
        self._rowMatching = [-1] * row
        self._colMatching = [-1] * col
        self._matchingEdges: Optional[List[Tuple[int, int]]] = None

    def addEdge(self, u: int, v: int) -> None:
        """男孩u和女孩v连边"""
        assert 0 <= u < self._row
        assert 0 <= v < self._col
        self._graph[u].append(v)

    def work(self) -> int:
        """返回最大匹配的个数"""

        def dfs(cur: int) -> bool:
            """寻找增广路"""
            if visited[cur]:
                return False
            visited[cur] = True
            for next in self._graph[cur]:
                if self._colMatching[next] == -1 or dfs(self._colMatching[next]):
                    self._colMatching[next] = cur
                    self._rowMatching[cur] = next
                    return True
            return False

        res = 0
        visited = [False] * self._row
        hasUpdated = True
        while hasUpdated:
            hasUpdated = False
            for cur in range(self._row):
                if self._rowMatching[cur] == -1 and dfs(cur):
                    hasUpdated = True
                    res += 1
            if hasUpdated:
                visited = [False] * self._row
        return res

    @property
    def matchingEdges(self):
        """返回匹配的边"""
        if self._matchingEdges is not None:
            return self._matchingEdges

        edges = []
        for cur in range(self._row):
            if self._rowMatching[cur] != -1:
                edges.append((cur, self._rowMatching[cur]))
        self._matchingEdges = edges
        return edges


###################################################################
# !deprecated
AdjList = Sequence[Iterable[int]]
AdjMap = Mapping[int, Iterable[int]]
Graph = Union[AdjList, AdjMap]  # 无向图


def hungarian(graph: Graph):
    """匈牙利算法求无权二分图最大匹配

    时间复杂度O(V * E)

    Args:
        graph (Graph): 无向图邻接表
    """

    def getColor(graph: Graph) -> Mapping[int, int]:
        """检测二分图并染色"""

        def dfs(cur: int, color: int) -> None:
            colors[cur] = color
            for next in graph[cur]:
                if colors[next] == -1:
                    dfs(next, color ^ 1)
                elif colors[cur] == colors[next]:
                    raise Exception("不是二分图")

        colors = defaultdict(lambda: -1)
        for i in range(n):
            if colors[i] == -1:
                dfs(i, 0)
        return colors

    def dfs(boy: int) -> bool:
        """寻找增广路"""
        nonlocal visited
        if boy in visited:
            return False
        visited.add(boy)

        for girl in graph[boy]:
            if matching[girl] == -1 or dfs(matching[girl]):
                matching[boy] = girl
                matching[girl] = boy
                return True
        return False

    n = len(graph)
    maxMatching = 0
    matching = defaultdict(lambda: -1)
    colors = getColor(graph)
    visited = set()
    for i in range(n):
        visited = set()
        if colors[i] == 0 and matching[i] == -1:
            if dfs(i):
                maxMatching += 1

    return maxMatching, matching, colors
