from typing import List, Set


def hungarian(adjList: List[Set[int]]) -> int:
    def getColor(adjList: List[Set[int]]) -> List[int]:
        """检测二分图并染色"""

        def dfs(cur: int, color: int) -> None:
            colors[cur] = color
            for next in adjList[cur]:
                if colors[next] == -1:
                    dfs(next, color ^ 1)
                elif colors[cur] == colors[next]:
                    raise Exception('不是二分图')

        n = len(adjList)
        colors = [-1] * n
        for i in range(n):
            if colors[i] == -1:
                dfs(i, 0)
        return colors

    def dfs(boy: int) -> bool:
        """寻找增广路"""
        nonlocal visited
        if visited[boy]:
            return False
        visited[boy] = True

        for girl in adjList[boy]:
            if matching[girl] == -1 or dfs(matching[girl]):
                matching[boy] = girl
                matching[girl] = boy
                return True
        return False

    n = len(adjList)
    maxMatching = 0
    matching = [-1] * n
    colors = getColor(adjList)
    visited = [False] * n
    for i in range(n):
        visited = [False] * n
        if colors[i] == 0 and matching[i] == -1:
            if dfs(i):
                maxMatching += 1

    return maxMatching
