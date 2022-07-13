# https://www.796t.com/article.php?id=594479

# !ProjectSelectionProblem
# !ProjectSelectionProblemは、「Xを選択してYを選択しない場合のペナルティ」を列挙していくのが基本である。

# !燃やす埋める問題に帰着
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

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")


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
    adjMap = defaultdict(lambda: defaultdict(int))
    for r in range(ROW):
        if rowSum[r] > 0:
            adjMap[PLAYER_A][r] += rowSum[r]
            res += rowSum[r]

    for c in range(COL):
        if colSum[c] > 0:
            adjMap[c + OFFSET][PLAYER_B] += colSum[c]
            res += colSum[c]

    for r in range(ROW):
        for c in range(COL):
            # 每個格子都可以通過割掉行列對應的點連到源點匯點的邊來成為只有行或者列選擇了它的格子。
            # 也可以割掉自己的邊來去掉重複貢獻。
            # !无穷表示这条边割不了
            adjMap[r][c + OFFSET] += Dinic.INF if matrix[r][c] < 0 else matrix[r][c]

    maxFlow = Dinic(adjMap)
    minCut = maxFlow.calMaxFlow(PLAYER_A, PLAYER_B)  # 每个点只能被1个玩家选择
    print(res - minCut)


if __name__ == "__main__":

    import sys
    from typing import DefaultDict
    from collections import defaultdict, deque

    Graph = DefaultDict[int, DefaultDict[int, int]]  # 有向带权图,权值为容量

    class Dinic:
        INF = int(1e18)

        def __init__(self, graph: Graph) -> None:
            self._graph = graph

        def calMaxFlow(self, start: int, end: int) -> int:
            def bfs() -> None:
                nonlocal depth, curArc
                depth = defaultdict(lambda: -1, {start: 0})
                visted = set([start])
                queue = deque([start])
                curArc = {cur: iter(self._reGraph[cur].keys()) for cur in self._reGraph.keys()}
                while queue:
                    cur = queue.popleft()
                    for child in self._reGraph[cur]:
                        if (child not in visted) and (self._reGraph[cur][child] > 0):
                            visted.add(child)
                            depth[child] = depth[cur] + 1
                            queue.append(child)

            def dfsWithCurArc(cur: int, minFlow: int) -> int:
                if cur == end:
                    return minFlow
                flow = 0
                while True:
                    if flow >= minFlow:
                        break

                    child = next(curArc[cur], None)
                    if child is not None:
                        if (depth[child] == depth[cur] + 1) and (self._reGraph[cur][child] > 0):
                            nextFlow = dfsWithCurArc(
                                child, min(minFlow - flow, self._reGraph[cur][child])
                            )
                            if nextFlow == 0:
                                depth[child] = -1
                            self._reGraph[cur][child] -= nextFlow
                            self._reGraph[child][cur] += nextFlow
                            flow += nextFlow
                    else:
                        break
                return flow

            self._updateRedisualGraph()

            res = 0
            depth = defaultdict(lambda: -1, {start: 0})
            curArc = dict()

            while True:
                bfs()
                if depth[end] != -1:
                    while True:
                        delta = dfsWithCurArc(start, Dinic.INF)
                        if delta == 0:
                            break
                        res += delta
                else:
                    break
            return res

        def getFlowOfEdge(self, v1: int, v2: int) -> int:
            """边的流量=容量-残量"""
            assert v1 in self._graph and v2 in self._graph[v1]
            return self._graph[v1][v2] - self._reGraph[v1][v2]

        def getRemainOfEdge(self, v1: int, v2: int) -> int:
            """边的残量(剩余的容量)"""
            assert v1 in self._graph and v2 in self._graph[v1]
            return self._reGraph[v1][v2]

        def _updateRedisualGraph(self) -> None:
            """残量图 存储每条边的剩余流量"""
            self._reGraph = defaultdict(lambda: defaultdict(int))
            for cur in self._graph:
                for next in self._graph[cur]:
                    self._reGraph[cur][next] = self._graph[cur][next]
                    self._reGraph[next].setdefault(cur, 0)

    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
