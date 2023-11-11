# 给一棵树，对每一个节点染成黑色或白色。

# 对于每一个节点，求强制把这个节点染成黑色的情况下，
# !所有的黑色节点组成一个联通块的染色方案数，答案对 M 取模。
# n<=1e5

# https://www.luogu.com.cn/problem/solution/AT4543

"""メモ 20210521
参照
https://algo-logic.info/tree-dp/
https://ei1333.hateblo.jp/entry/2017/04/10/224413
https://qiita.com/keymoon/items/2a52f1b0fb7ef67fb89e
問題
https://atcoder.jp/contests/dp/tasks/dp_v
"""


from Rerooting import Rerooting

if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    INF = int(4e18)

    def e(root: int) -> int:
        # 根节点的染色方案为1
        return 1

    def op(childRes1: int, childRes2: int) -> int:
        # 子树的染色方案就是他所有子树染色方案 +1 的积
        return childRes1 * childRes2 % MOD

    def composition(fromRes: int, parent: int, cur: int, direction: int) -> int:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        # 子树的染色方案就是他所有子树染色方案 +1 的积
        # +1表示子树可以不染色
        return fromRes + 1

    n, MOD = map(int, input().split())
    R = Rerooting(n, decrement=0)  # 0-indexed
    for _ in range(n - 1):
        x, y = map(int, input().split())
        x, y = x - 1, y - 1
        R.addEdge(x, y)

    dp = R.rerooting(composition=composition, op=op, e=e, root=0)
    print(*dp, sep="\n")
