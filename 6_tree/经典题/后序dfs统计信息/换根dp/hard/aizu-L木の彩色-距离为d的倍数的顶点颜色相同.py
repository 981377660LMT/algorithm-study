# https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=3148&lang=ja
# 树的涂色方案数
# 对每个d (1<=d<=n) 问: 是否能将这棵树涂色,
# !满足:距离为d的倍数的顶点颜色相同,距离不为d的倍数的顶点颜色不同


from Rerooting import Rerooting
from typing import Tuple

INF = int(1e18)
if __name__ == "__main__":
    E = Tuple[int, int, int]  # 子树内不同分支的最大三个长度。

    def e(root: int) -> E:
        return (-INF, -INF, -INF)

    def op(childRes1: E, childRes2: E) -> E:
        """合并当前根的不同分支长度,最多保留3个"""
        a1, b1, c1 = childRes1
        a2, b2, c2 = childRes2
        tmp = [a1, b1, c1, a2, b2, c2]
        tmp.sort(reverse=True)
        return (tmp[0], tmp[1], tmp[2])

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        """上传时只保留子树内的最大深度"""
        a = fromRes[0]
        a = 0 if a < 0 else a
        a += 1  # 边权？？
        return (a, -INF, -INF)

    n = int(input())
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        edges.append((u - 1, v - 1))
    if n == 1:
        print("1")
        exit(0)
    R = Rerooting(n)
    for u, v in edges:
        R.addEdge(u, v)
    tree = R.adjList

    # dp上只剩下每个点作为根时，不同分支的最大三个长度。
    dp = R.rerooting(e=e, op=op, composition=composition, root=0)
    res = [True] * (n + 1)
    for i in range(n):
        if len(tree[i]) < 3:
            continue
        a, b, c = dp[i]
        if a == b == c:
            # !那么a到b的倒数第二格，a到c的倒数第二格颜色必须相同，但b的倒数第二格到c的倒数第二格颜色相同且距离不同，不符合条件；
            res[a + a - 1] = False
        else:  # a>c
            # !推出矛盾
            res[a + c] = False

    # 某个更大的长度无法达到，更小的长度也一定无法达到
    # 如果5不行 那么肯定不行 可以理解为把b和c里的点向上平移一层
    # 到顶了就移a
    for i in range(n - 1, -1, -1):
        if not res[i + 1]:
            res[i] = False

    res[1] = True
    res[2] = True
    res = res[1:]
    print("".join(map(str, map(int, res))))
