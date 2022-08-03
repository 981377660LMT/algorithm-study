# https://atcoder.jp/contests/abc241/tasks/abc241_g

# 巡回对战(総当たり戦/round robin )的橄榄球比赛
# n个队每个队要打n-1场比赛 每场比赛必须分出胜负
# 现在已知m场比赛的结果
# !对每个人x判断是否可能成为胜者(胜利场数最多且唯一)。
# n<=50

# https://www.bilibili.com/read/cv15437330?spm_id_from=333.1007.0.0

# 橄榄球比赛建模(Baseball elimination) => 最大流算法
# !对每个人，看他是否能够取胜 即要求除x外每个人的胜利场数小于定值
# !建图，有O(n^2)条边 容量为O(n) 所以Dinic每次查询复杂度不超过O(maxflow*E) 即 O(n^3)

# 源点 => 分配胜利给每个场次
# 每个场次 => 将胜利分给两队中的一个
# 每个队最多还可以赢几场比赛 => 流向汇点
# !如果最大流等于要分配的场数 那么就可能成为获胜者


# TODO 有问题 WA8 TLE1 AC63
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(1e18)


def main() -> None:
    def query(person: int) -> bool:
        """person是否有胜利的可能

        假设他之后全胜 别人必须比他的得分严格小
        对其他人的比赛进行最大流建模
        如果最大流等于要分配的场数 那么存在分配方案
        即person可能成为胜者
        """
        maxScore = win[person] + remain[person]  # 别人必须比这个小
        if maxScore <= maxWin:
            return False

        START, END, OFFSET = -1, -2, int(1e6)
        maxFlow = MaxFlow(START, END)

        remainWin = 0  # 还需分配的场数
        for i in range(n):
            if i == person:
                continue
            for j in range(i + 1, n):
                if j == person:
                    continue
                if todo[i][j] > 0:
                    remainWin += todo[i][j]
                    game = i * n + j
                    maxFlow.addEdge(START, game, todo[i][j])
                    maxFlow.addEdge(game, i + OFFSET, todo[i][j])
                    maxFlow.addEdge(game, j + OFFSET, todo[i][j])

        for i in range(n):
            if i == person:
                continue
            maxFlow.addEdge(i + OFFSET, END, maxScore - 1 - win[i])  # 别的队还可以最多赢几场比赛

        return maxFlow.calMaxFlow() == remainWin  # 最大流等于还需分配的场数

    n, m = map(int, input().split())

    win, lose, remain = [0] * n, [0] * n, [n - 1] * n
    todo = [[1] * n for _ in range(n)]  # !i与j之间剩余的对战数
    for i in range(n):
        todo[i][i] = 0

    for _ in range(m):
        u, v = map(int, input().split())  # !u赢了v
        u, v = u - 1, v - 1
        win[u] += 1
        lose[v] += 1
        remain[u] -= 1
        remain[v] -= 1
        todo[u][v] -= 1
        todo[v][u] -= 1

    maxWin = max(win)
    res = []
    for i in range(n):
        if query(i):
            res.append(i + 1)
    print(*res)


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

        def addEdge(self, v1: int, v2: int, w: int, *, cover=False) -> None:
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
