import collections

# 你需要找到起点到终点的最短路径，然而你不知道网格的大小、起点和终点。你只能向 GridMaster 对象查询。

# dfs建图 + Dijkstra求成本最小路径
class GridMaster(object):
    def canMove(self, direction: str) -> bool:
        ...

    def move(self, direction: str) -> int:
        ...

    def isTarget(self) -> None:
        ...


class Solution(object):
    def findShortestPath(self, master: 'GridMaster') -> int:

        valid = set()
        dirs = {'U': [-1, 0], 'D': [1, 0], 'L': [0, -1], 'R': [0, 1]}
        back = {'U': 'D', 'D': 'U', 'L': 'R', 'R': 'L'}
        target = None

        def dfs(x, y):
            nonlocal target

            if master.isTarget():
                target = (x, y)

            for d in dirs:
                if master.canMove(d):
                    dx, dy = dirs[d]
                    i, j = x + dx, y + dy
                    if (i, j) not in valid:
                        valid.add((x, y))
                        master.move(d)
                        dfs(i, j)
                        master.move(back[d])

        dfs(0, 0)

        if target == None:
            return -1

        visited = set()
        queue = collections.deque([(0, 0)])
        dis = 0
        while queue:
            for _ in range(len(queue)):
                x, y = queue.popleft()
                if (x, y) == target:
                    return dis
                for i, j in [[x + 1, y], [x - 1, y], [x, y + 1], [x, y - 1]]:
                    if (i, j) in valid and (i, j) not in visited:
                        visited.add((i, j))
                        queue.append((i, j))
            dis += 1

        return -1

