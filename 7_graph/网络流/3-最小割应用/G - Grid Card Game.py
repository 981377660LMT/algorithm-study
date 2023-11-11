# https://www.796t.com/article.php?id=594479

# !ProjectSelectionProblem
# !ProjectSelectionProblemは、「Xを選択してYを選択しない場合のペナルティ」を列挙していくのが基本である。

# !燃やす埋める問題に帰着 最小割
# 给定一个矩阵Anxn (1≤n ≤ 100)，选择一些行列，可以得到这些行列包含的位置的并的数值和。
# 此外要求任意选中的行列交点处不能是负数。
# !求选择的最大值
# !时间复杂度O(V^2E)

# 重新描述一下問題就是：
# 1.可以選擇一些行列，給答案加上行列的權值和。
# 2.一個負權格子不能同時選上它在的行列。
# 3.一個非負權格子行列同時被選上的時候要減掉自己的權值（因為它貢獻給答案貢獻了兩次）。

# 首先考慮先選擇所有對答案貢獻非負的行和列，然後就只要考慮如何去掉重複貢獻的非負權格子和強制不同時選擇負權格子所在的行列。
# 建圖，每一非負整行、每一非負整列都存在對應結點。從源點向每個貢獻非負的行點連最大流量為行權值和的邊，列同理，但聯向匯點。
# 每個格子就在所在行列之間連邊，負點就連無窮大的權值，代表不能割這條邊，非負點就連自己格子的權值，割掉就是刪去自己重複的貢獻。
# !這樣每個格子都可以通過割掉行列對應的點連到源點匯點的邊來成為只有行或者列選擇了它的格子。也可以割掉自己的邊來去掉重複貢獻。

import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
INF = int(1e18)


def main() -> None:
    ROW, COL = map(int, input().split())
    matrix = [tuple(map(int, input().split())) for _ in range(ROW)]
    rowSum = [0] * ROW
    colSum = [0] * COL

    for r in range(ROW):
        for c in range(COL):
            rowSum[r] += matrix[r][c]
            colSum[c] += matrix[r][c]

    res = 0
    PLAYER_A, PLAYER_B, OFFSET = -1, -2, int(1e4)
    maxFlow = MaxFlow(PLAYER_A, PLAYER_B)

    for r in range(ROW):
        if rowSum[r] > 0:
            maxFlow.addEdge(PLAYER_A, r, rowSum[r])
            res += rowSum[r]

    for c in range(COL):
        if colSum[c] > 0:
            maxFlow.addEdge(c + OFFSET, PLAYER_B, colSum[c])
            res += colSum[c]

    for r in range(ROW):
        for c in range(COL):
            # 每個格子都可以通過割掉行列對應的點連到源點匯點的邊來成為只有行或者列選擇了它的格子。
            # 也可以割掉自己的邊來去掉重複貢獻。
            # !无穷表示这条边割不了
            if matrix[r][c] < 0:
                maxFlow.addEdge(r, c + OFFSET, INF)
            else:
                maxFlow.addEdge(r, c + OFFSET, matrix[r][c])

    minCut = maxFlow.calMaxFlow()  # 每个点只能被1个玩家选择
    print(res - minCut)


if __name__ == "__main__":

    from collections import defaultdict, deque
    from typing import Set

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

    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
