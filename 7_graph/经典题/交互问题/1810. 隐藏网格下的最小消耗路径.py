import collections
from heapq import heappop, heappush

# 这是一个交互问题。
# 你需要找到机器人移动至目标网格的最小总消耗。但可惜的是你并不知道网格的尺寸、初始单元和目标单元。你只允许通过询问GridMaster类获得信息。

# dfs建图 + Dijkstra求成本最小路径


class GridMaster(object):
    def canMove(self, direction: str) -> bool:
        ...

    # 沿该方向移动机器人，并返回移动到该单元的消耗值(边的weight)
    def move(self, direction: str) -> int:
        ...

    def isTarget(self) -> None:
        ...


class Solution(object):
    def findShortestPath(self, master: 'GridMaster') -> int:
        valid = dict({(0, 0): 0})
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
                        cost = master.move(d)
                        valid[(i, j)] = cost
                        dfs(i, j)
                        master.move(back[d])

        dfs(0, 0)

        if target == None:
            return -1

        # 此处dist无法用数组,因为不知长度
        dist = dict({(0, 0): 0})
        pq = [(0, (0, 0))]
        while pq:
            cost, (x, y) = heappop(pq)
            for i, j in [[x - 1, y], [x + 1, y], [x, y - 1], [x, y + 1]]:
                if (i, j) in valid and cost + valid[(i, j)] < dist.get((i, j), 0x7FFFFFFF):
                    heappush(pq, (cost + valid[(i, j)], (i, j)))
                    dist[(i, j)] = cost + valid[(i, j)]

        return dist[target]


# 输入: grid = [[2,3],[1,1]], r1 = 0, c1 = 1, r2 = 1, c2 = 0
# 输出: 2
# 解释: 其中一种可能路径描述如下：
# 机器人最开始站在单元格 (0, 1) ，用 3 表示
# - master.canMove('U') 返回 false
# - master.canMove('D') 返回 true
# - master.canMove('L') 返回 true
# - master.canMove('R') 返回 false
# - master.move('L') 机器人移动到单元格 (0, 0) 并返回 2
# - master.isTarget() 返回 false
# - master.canMove('U') 返回 false
# - master.canMove('D') 返回 true
# - master.canMove('L') 返回 false
# - master.canMove('R') 返回 true
# - master.move('D') 机器人移动到单元格 (1, 0) 并返回 1
# - master.isTarget() 返回 true
# - master.move('L') 机器人不移动并返回 -1
# - master.move('R') 机器人移动到单元格 (1, 1) 并返回 1
# 现在我们知道了机器人达到目标单元(1, 0)的最小消耗成本为2。

