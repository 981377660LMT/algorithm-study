# https://atcoder.jp/contests/practice2/submissions/32497027
# 给定 n 个还未赋值的布尔变量 x1∼xn。
# 现在有 m 个条件，每个条件的形式为
#  !“xi 为 0/1 或 xj 为 0/1 至少有一项成立”，
# !例如 “x1 为 1 或 x3 为 0”、“x8 为 0 或 x4 为 0” 等。
# 现在，请你对这 n 个布尔变量进行赋值（0 或 1），使得所有 m 个条件能够成立。


from typing import List


class TwoSAT:
    def __init__(self, n: int):
        """n个条件,每个条件的形式都是 "xi 为 true/false 或 xj 为 true/false 至少有一项成立"

        2-SAT 问题的目标是给每个变量赋值，使得所有条件成立
        """
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

    def addEdge(self, i: int, iState: bool, j: int, jState: bool) -> None:
        """对于 Xi 为真，可以通过连边 <i+n,i> 实现; Xi 为假，可以通过连边 <i,i+n> 来实现"""
        i0 = i
        i1 = i + self._n
        if not iState:
            i0, i1 = i1, i0
        j0 = j
        j1 = j + self._n
        if not jState:
            j0, j1 = j1, j0
        self._adjList[i1].append(j0)
        self._adjList[j1].append(i0)

    def buildGraph(self) -> None:
        _, self.group = self._getSCC(2 * self._n, self._adjList)

    def check(self) -> bool:
        """强连通分量中同时存在 a => a非 及 a非 => a 时，命题无解"""
        for i in range(self._n):
            if self.group[i] == self.group[i + self._n]:
                return False
        return True

    def work(self) -> List[bool]:
        """每个条件结点的正确性"""
        res = [False] * self._n
        for i in range(self._n):
            if self.group[i] > self.group[i + self._n]:  # type: ignore
                res[i] = True
        return res


if __name__ == "__main__":
    n, m = map(int, input().split())
    twoSAT = TwoSAT(n)
    for _ in range(m):
        i, iState, j, jState = map(int, input().split())
        twoSAT.addEdge(i - 1, bool(iState), j - 1, bool(jState))

    twoSAT.buildGraph()

    if not twoSAT.check():
        print("IMPOSSIBLE")
        exit(0)

    print("POSSIBLE")
    res = twoSAT.work()
    print(*[int(b) for b in res])
