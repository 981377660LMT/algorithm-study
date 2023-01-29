# "." 表示恶魔可随意通行的平地；
# "#" 表示恶魔不可通过的障碍物，玩家可通过在 `平地` 上设置障碍物，即将 "." 变为 "#" 以阻挡恶魔前进；
# !(注意只能割平地)
# "S" 表示恶魔出生点，将有大量的恶魔该点生成，恶魔可向上/向下/向左/向右移动，且无法移动至地图外；
# "P" 表示瞬移点，移动到 "P" 点的恶魔可被传送至任意一个 "P" 点，也可选择不传送；
# "C" 表示城堡。
# 地图上可能有一个或多个出生点
# 地图上有且只有一个城堡

# !我们希望「恶魔」无法到达城堡，即「城堡」和「恶魔」之间不连通，对应到图论模型上就是「割」问题。
# !建立源点向恶魔出生点连边，再将城堡向汇点连边；再将图中相邻的点之间连边即可
# 传送门互相连通，建立一个特殊点把它们都免费连起来即可。
# 因为要移除点而不是移除边 所以要把所有点拆成in 和 out
# "."的自身流量为 1，其余边的流量为 inf,


# https://leetcode.cn/problems/7rLGCR/solution/lcp-38-shou-wei-cheng-bao-by-zerotrac2-kgv2/
# https://leetcode.cn/problems/7rLGCR/solution/czui-xiao-ge-jian-mo-dai-zhu-shi-by-litt-g77h/


# !这样最小割一定发生在空地上，且最小割的大小就是建立障碍的个数，也就是答案
from collections import defaultdict, deque
from typing import Generator, List, Set


DIR4 = [[0, 1], [0, -1], [1, 0], [-1, 0]]
INF = int(1e18)


class Solution:
    def guardCastle(self, grid: List[str]) -> int:
        def genNext(point: int) -> Generator[int, None, None]:
            curRow, curCol = point // COL, point % COL
            for dRow, dCol in DIR4:
                nextRow, nextCol = curRow + dRow, curCol + dCol
                if 0 <= nextRow < ROW and 0 <= nextCol < COL and grid[nextRow][nextCol] != "#":
                    yield nextRow * COL + nextCol

        ROW, COL = len(grid), len(grid[0])
        START, END, OFFSET, TELEPORT = -1, -2, int(1e5), -3
        mf = MaxFlow(START, END)

        for r in range(ROW):
            for c in range(COL):
                if grid[r][c] == "#":
                    continue
                cur = r * COL + c

                # 0. 所有点拆成 入点 和 出点 两个点
                mf.addEdge(cur, cur + OFFSET, 1 if grid[r][c] == "." else INF)

                # 1. 源点连接恶魔出生点
                if grid[r][c] == "S":
                    mf.addEdge(START, cur, INF)

                # 2. 城堡连接汇点
                if grid[r][c] == "C":
                    mf.addEdge(cur, END, INF)

                # 3. 虚拟点连通所有传送门
                if grid[r][c] == "P":
                    mf.addEdge(cur + OFFSET, TELEPORT, INF)
                    mf.addEdge(TELEPORT, cur, INF)

                # 4. 所有出点连通周围的入点
                for next in genNext(cur):
                    mf.addEdge(cur + OFFSET, next, INF)

        minCut = mf.calMaxFlow()
        return minCut if minCut < INF else -1


class MaxFlow:
    def __init__(self, start: int, end: int) -> None:
        self.graph = defaultdict(lambda: defaultdict(int))  # 原图
        self._start = start
        self._end = end

    def calMaxFlow(self) -> int:
        self._updateRedisualGraph()
        start, end = self._start, self._end
        flow = 0

        while self._bfs():
            delta = INF
            while delta:
                delta = self._dfs(start, end, INF)
                flow += delta
        return flow

    def addEdge(self, v1: int, v2: int, w: int, *, cover=True) -> None:
        """添加边 v1->v2, 容量为w

        Args:
            v1: 边的起点
            v2: 边的终点
            w: 边的容量
            cover: 是否覆盖原有边
        """
        if cover:
            self.graph[v1][v2] = w
        else:
            self.graph[v1][v2] += w

    def getFlowOfEdge(self, v1: int, v2: int) -> int:
        """边的流量=容量-残量"""
        assert v1 in self.graph and v2 in self.graph[v1]
        return self.graph[v1][v2] - self._reGraph[v1][v2]

    def getRemainOfEdge(self, v1: int, v2: int) -> int:
        """边的残量(剩余的容量)"""
        assert v1 in self.graph and v2 in self.graph[v1]
        return self._reGraph[v1][v2]

    def getPath(self) -> Set[int]:
        """最大流经过了哪些点"""
        visited = set()
        stack = [self._start]
        reGraph = self._reGraph
        while stack:
            cur = stack.pop()
            visited.add(cur)
            for next, remain in reGraph[cur].items():
                if next not in visited and remain > 0:
                    visited.add(next)
                    stack.append(next)
        return visited

    def _updateRedisualGraph(self) -> None:
        """残量图 存储每条边的剩余流量"""
        self._reGraph = defaultdict(lambda: defaultdict(int))
        for cur in self.graph:
            for next, cap in self.graph[cur].items():
                self._reGraph[cur][next] = cap
                self._reGraph[next].setdefault(cur, 0)  # 注意自环边

    def _bfs(self) -> bool:
        self._depth = depth = defaultdict(lambda: -1, {self._start: 0})
        reGraph, start, end = self._reGraph, self._start, self._end
        queue = deque([start])
        self._iters = {cur: iter(reGraph[cur].keys()) for cur in reGraph.keys()}
        while queue:
            cur = queue.popleft()
            nextDist = depth[cur] + 1
            for next, remain in reGraph[cur].items():
                if depth[next] == -1 and remain > 0:
                    depth[next] = nextDist
                    queue.append(next)

        return depth[end] != -1

    def _dfs(self, cur: int, end: int, flow: int) -> int:
        if cur == end:
            return flow
        reGraph, depth, iters = self._reGraph, self._depth, self._iters
        for next in iters[cur]:
            remain = reGraph[cur][next]
            if remain and depth[cur] < depth[next]:
                nextFlow = self._dfs(next, end, min(flow, remain))
                if nextFlow:
                    reGraph[cur][next] -= nextFlow
                    reGraph[next][cur] += nextFlow
                    return nextFlow
        return 0


print(Solution().guardCastle(grid=["S.C.P#P.", ".....#.S"]))
print(Solution().guardCastle(grid=["SP#P..P#PC#.S", "..#P..P####.#"]))
print(Solution().guardCastle(grid=["SP#.C.#PS", "P.#...#.P"]))
print(Solution().guardCastle(grid=["CP.#.P.", "...S..S"]))
