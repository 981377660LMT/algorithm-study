# https://yukicoder.me/problems/no/1075
# yukicoder No.1075-成为山的方案数
# 每个顶点分配[1,k]中的任一个数
# !使得存在一个顶点i,到其他任意点j的路径都是单调不减的
# 求顶点分配方案数(山的个数)模1e9+7


# !怎么去重?
# 只计算离父亲最近的点

from Rerooting import Rerooting


from itertools import accumulate
from typing import List


INF = int(4e18)
MOD = int(1e9 + 7)

if __name__ == "__main__":

    E = List[int]  # 子树中每种取值的方案数

    def e(root: int) -> E:
        return [1] * k

    def op(e1: E, e2: E) -> E:
        return [e1[i] * e2[i] % MOD for i in range(k)]

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        return list(accumulate(fromRes, lambda x, y: x + y % MOD))

    n, k = map(int, input().split())
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        edges.append((u - 1, v - 1))

    R = Rerooting(n)
    for u, v in edges:
        R.addEdge(u, v)

    tree = R.adjList

    # !dp[i][j]: 每个顶点i作为山峰时,分配j(1<=j<=k)有多少种方案
    dp = R.rerooting(e=e, op=op, composition=composition, root=0)

    res = 0

    # !不对
    def dfs(cur: int, pre: int) -> None:
        global res
        sub = [1] * k
        for next in tree[cur]:
            if next == pre:
                sub[0] = 0
                for i in range(1, k):
                    sub[i] *= dp[next][i - 1]
                    sub[i] %= MOD
            else:
                for i in range(k):
                    sub[i] *= dp[next][i]
                    sub[i] %= MOD
                dfs(next, cur)
        for v in sub:
            res += v
            res %= MOD

    dfs(0, -1)
    print(res)
