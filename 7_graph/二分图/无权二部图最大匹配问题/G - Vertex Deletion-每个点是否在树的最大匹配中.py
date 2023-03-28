# abc-233-G - Vertex Deletion-每个点是否在树的最大匹配中
# https://atcoder.jp/contests/abc223/tasks/abc223_g
# 给定一棵树
# 对每个结点i为根,删除根连接的所有边后,
# !使得剩下的树的最大匹配和原树最大匹配相等
# 求这样的根的个数

# !解:即不参与二分图的最大匹配

# 类似:
# https://yukicoder.me/problems/2085
# 二分图博弈
# Alice和Bob在树上博弈
# 先手放一个棋子,后手在相邻的结点放一个棋子
# 交替放棋子,直到不能放棋子的时候,输
# !问先手是否必胜 => 如果起点不在二分图的最大匹配中,先手必胜

from typing import Callable, Generic, List, TypeVar

T = TypeVar("T")


class Rerooting(Generic[T]):

    __slots__ = ("adjList", "_n", "_decrement")

    def __init__(self, n: int, decrement: int = 0):
        self.adjList = [[] for _ in range(n)]
        self._n = n
        self._decrement = decrement

    def addEdge(self, u: int, v: int) -> None:
        u -= self._decrement
        v -= self._decrement
        self.adjList[u].append(v)
        self.adjList[v].append(u)

    def rerooting(
        self,
        e: Callable[[int], T],
        op: Callable[[T, T], T],
        composition: Callable[[T, int, int, int], T],
        root=0,
    ) -> List["T"]:
        root -= self._decrement
        assert 0 <= root < self._n
        parents = [-1] * self._n
        order = [root]
        stack = [root]
        while stack:
            cur = stack.pop()
            for next in self.adjList[cur]:
                if next == parents[cur]:
                    continue
                parents[next] = cur
                order.append(next)
                stack.append(next)

        dp1 = [e(i) for i in range(self._n)]
        dp2 = [e(i) for i in range(self._n)]
        for cur in order[::-1]:
            res = e(cur)
            for next in self.adjList[cur]:
                if parents[cur] == next:
                    continue
                dp2[next] = res
                res = op(res, composition(dp1[next], cur, next, 0))
            res = e(cur)
            for next in self.adjList[cur][::-1]:
                if parents[cur] == next:
                    continue
                dp2[next] = op(res, dp2[next])
                res = op(res, composition(dp1[next], cur, next, 0))
            dp1[cur] = res

        for newRoot in order[1:]:
            parent = parents[newRoot]
            dp2[newRoot] = composition(op(dp2[newRoot], dp2[parent]), parent, newRoot, 1)
            dp1[newRoot] = op(dp1[newRoot], dp2[newRoot])
        return dp1


if __name__ == "__main__":
    n = int(input())
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        edges.append((u - 1, v - 1))

    E = int  # 当前节点是否构成子树的最大匹配, 0: 不参与, 1: 参与

    def e(root: int) -> E:
        return 0

    def op(childRes1: E, childRes2: E) -> E:
        return childRes1 | childRes2

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        return fromRes ^ 1  # 孩子参与匹配则父亲不参与, 反之成立

    R = Rerooting(n)
    for u, v in edges:
        R.addEdge(u, v)

    dp = R.rerooting(e=e, op=op, composition=composition, root=0)

    print(dp.count(0))  # 不在最大匹配中的点的个数
