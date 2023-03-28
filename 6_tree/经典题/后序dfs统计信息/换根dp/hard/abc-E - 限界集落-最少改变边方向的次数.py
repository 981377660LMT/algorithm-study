# https://atcoder.jp/contests/njpc2017/tasks/njpc2017_e
# 给定一棵树和边,每条边有一个好方向和坏方向, 沿着坏方向走有罚款.
# !求一个最好的点，满足到最远结点距离在d之内, 并且到这个点的罚款数最少, 求出最少的罚款数.
# 不存在输出-1


from Rerooting import Rerooting

from collections import defaultdict
from typing import Tuple


INF = int(4e18)

if __name__ == "__main__":

    E = Tuple[int, int]  # (最远距离，总罚款数)

    def e(root: int) -> E:
        return (0, 0)

    def op(childRes1: E, childRes2: E) -> E:
        dist1, pay1 = childRes1
        dist2, pay2 = childRes2
        return (max(dist1, dist2), pay1 + pay2)

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        from_ = parent if direction == 1 else cur
        to = cur if direction == 1 else parent
        dist, pay = fromRes
        return (dist + weight[from_][to], pay + penalty[from_][to])

    n, d = map(int, input().split())
    R = Rerooting(n)
    penalty = [defaultdict(int) for _ in range(n)]
    weight = [defaultdict(int) for _ in range(n)]
    for _ in range(n - 1):
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        R.addEdge(u, v)
        weight[u][v] = w
        weight[v][u] = w
        penalty[v][u] = 1  # 反向走代价为1

    dp = R.rerooting(e=e, op=op, composition=composition, root=0)
    ok = [i for i in range(n) if dp[i][0] <= d]
    if not ok:
        print(-1)
        exit(0)
    print(min(dp[i][1] for i in ok))
