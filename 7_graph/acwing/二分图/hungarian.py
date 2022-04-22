from collections import defaultdict
from typing import DefaultDict, Set


def hungarian(adjMap: DefaultDict[int, Set[int]]) -> int:
    def getColor(adjMap: DefaultDict[int, Set[int]]) -> DefaultDict[int, int]:
        """检测二分图并染色"""

        def dfs(cur: int, color: int) -> None:
            colors[cur] = color
            for next in adjMap[cur]:
                if colors[next] == -1:
                    dfs(next, color ^ 1)
                elif colors[cur] == colors[next]:
                    raise Exception('不是二分图')

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

        for girl in adjMap[boy]:
            if matching[girl] == -1 or dfs(matching[girl]):
                matching[boy] = girl
                matching[girl] = boy
                return True
        return False

    n = len(adjMap)
    maxMatching = 0
    matching = defaultdict(lambda: -1)
    colors = getColor(adjMap)
    visited = set()
    for i in range(n):
        visited = set()
        if colors[i] == 0 and matching[i] == -1:
            if dfs(i):
                maxMatching += 1

    return maxMatching

