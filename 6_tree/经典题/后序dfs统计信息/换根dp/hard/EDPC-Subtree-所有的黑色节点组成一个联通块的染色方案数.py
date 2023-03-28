# https://blog.hamayanhamayan.com/entry/2019/01/12/152428
# https://atcoder.jp/contests/dp/tasks/dp_v

# 给一棵树，对每一个节点染成黑色或白色。
# 对于每一个节点，求强制把这个节点染成黑色的情况下，
# !所有的黑色节点组成一个联通块的染色方案数，答案对 MOD 取模。


import sys
from Rerooting import Rerooting

input = lambda: sys.stdin.readline().rstrip("\r\n")
INF = int(4e18)

if __name__ == "__main__":

    E = int

    def e(root: int) -> E:
        return 1

    def op(childRes1: E, childRes2: E) -> E:
        return childRes1 * childRes2 % MOD

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        return fromRes + 1  # 加一是整棵子树全为白色的方案数

    n, MOD = map(int, input().split())
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        edges.append((u - 1, v - 1))

    R = Rerooting(n)
    for u, v in edges:
        R.addEdge(u, v)

    dp = R.rerooting(e=e, op=op, composition=composition, root=0)

    print(*dp, sep="\n")
