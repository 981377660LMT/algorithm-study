# https://www.luogu.com.cn/problem/P4782
# 给定 n 个还未赋值的布尔变量 x1∼xn。
# 现在有 m 个条件，每个条件的形式为
#  !“xi 为 0/1 或 xj 为 0/1 至少有一项成立”，
# !例如 “x1 为 1 或 x3 为 0”、“x8 为 0 或 x4 为 0” 等。
# 现在，请你对这 n 个布尔变量进行赋值（0 或 1），使得所有 m 个条件能够成立。


from typing import List


class TwoSat:
    def __init__(self, n: int):
        self._n = n
        self._adjList = [[] for _ in range(2 * n)]

    # a V B
    # pos_i = True -> a = i
    # pos_i = False -> a = ¬i
    @staticmethod
    def _getSCC(n: int, adjList: List[List[int]]):
        order = []
        visited = [False] * n
        group = [None] * n
        rAdjList = [[] for _ in range(n)]
        for i in range(n):
            for j in adjList[i]:
                rAdjList[j].append(i)

        def dfs(cur: int) -> None:
            stack = [(1, cur), (0, cur)]
            while stack:
                t, cur = stack.pop()
                if t == 0:
                    if visited[cur]:
                        stack.pop()
                        continue
                    visited[cur] = True
                    for npos in adjList[cur]:
                        if not visited[npos]:
                            stack.append((1, npos))
                            stack.append((0, npos))
                else:
                    order.append(cur)

        def rdfs(cur: int, dfsId: int) -> None:
            stack = [cur]
            group[cur] = dfsId
            visited[cur] = True
            while stack:
                cur = stack.pop()
                for npos in rAdjList[cur]:
                    if not visited[npos]:
                        visited[npos] = True
                        group[npos] = dfsId
                        stack.append(npos)

        for i in range(n):
            if not visited[i]:
                dfs(i)

        visited = [False] * n
        count = 0
        for s in reversed(order):
            if not visited[s]:
                rdfs(s, count)
                count += 1
        return count, group

    def addLimit(self, i: int, iState: bool, j: int, jState: bool) -> None:
        """加边方式1:根据限制条件 '至少满足一个' 添加边.

        !i为iState 和 j为jState 两个条件至少满足一个(state 表示 真/假)
        则 `否i => j` 和 `否j => i`

        0 <= i < n, 0 <= j < n
        """
        notU, v = i + iState * self._n, j + (jState ^ 1) * self._n
        self._adjList[notU].append(v)
        notV, u = j + jState * self._n, i + (iState ^ 1) * self._n
        self._adjList[notV].append(u)

    def addEdge(self, u: int, v: int) -> None:
        """加边方式2:根据命题的推导关系添加边.

        如果 u => v (u成立可以推导出v成立)，那么就添加边 u => v 以及 ¬v => ¬u.。
        注意当 u/v 为真命题时, 0 <= u/v < n
        当 u/v 为假命题时, n <= u/v < 2*n
        """
        self._adjList[u].append(v)
        notU = (u + self._n) if u < self._n else (u - self._n)
        notV = (v + self._n) if v < self._n else (v - self._n)
        self._adjList[notV].append(notU)

    def build(self) -> None:
        _, self.group = self._getSCC(2 * self._n, self._adjList)

    def check(self) -> bool:
        """强连通分量中同时存在 a => 非a 及 非a => a 时，命题无解"""
        for i in range(self._n):
            if self.group[i] == self.group[i + self._n]:
                return False
        return True

    def work(self) -> List[bool]:
        """每个命题xi是否为真"""
        res = [False] * self._n
        for i in range(self._n):
            # !DAG上 如果 真命题的拓扑序排在假命题后面，说明可以推导出真命题，即为真
            res[i] = self.group[i] > self.group[i + self._n]  # type: ignore
        return res


if __name__ == "__main__":
    n, m = map(int, input().split())
    twoSAT = TwoSat(n)
    for _ in range(m):
        i, iState, j, jState = map(int, input().split())
        twoSAT.addLimit(i - 1, bool(iState), j - 1, bool(jState))
        # 两个条件至少满足一个

    twoSAT.build()

    if not twoSAT.check():
        print("IMPOSSIBLE")
        exit(0)

    print("POSSIBLE")
    res = twoSAT.work()
    print(*[int(b) for b in res])
