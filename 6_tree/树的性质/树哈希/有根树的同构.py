# https://judge.yosupo.jp/submission/102920
# rngハッシュ

# !每个点作为根节点时,求出树的哈希值,然后判断哈希值是否相同即可
# 只要求形态相同即可,不要求每个点的值也相同


from random import randint
from typing import List, Tuple


def rootedTreeIsomorphism(
    n: int, edges: List[Tuple[int, int]], root=0, mod=(1 << 61) - 1
) -> List[int]:
    tree = [[] for _ in range(n)]
    for u, cur in edges:
        tree[u].append(cur)
        tree[cur].append(u)

    parent = [-1] * n
    height = [0] * n
    order = [root]
    stack = [root]
    while stack:
        cur = stack.pop()
        for next in tree[cur]:
            if parent[cur] == next:
                continue
            parent[next] = cur
            order.append(next)
            stack.append(next)

    for i in range(len(order) - 1, 0, -1):
        cur = order[i]
        cand = height[cur] + 1
        if cand > height[parent[cur]]:
            height[parent[cur]] = cand

    dp = [1] * n
    rands = [randint(0, mod - 1) for _ in range(n)]
    for i in range(len(order) - 1, -1, -1):
        cur = order[i]
        h = height[cur]
        r = rands[h]
        for next in tree[cur]:
            if parent[cur] == next:
                continue
            dp[cur] *= r + dp[next]
            dp[cur] %= mod
    return dp


if __name__ == "__main__":
    n = int(input())
    parents = list(map(int, input().split()))

    edges = []
    for cur, pre in enumerate(parents):
        edges.append((cur + 1, pre))

    # 以每个根作为根节点求出树的哈希值
    dp = rootedTreeIsomorphism(n, edges)
    allNums = sorted(set(dp))
    mp = {v: k for k, v in enumerate(allNums)}
    res = [mp[h] for h in dp]

    print(len(allNums))  # 哈希值的种类数
    print(*res)  # 每个根的哈希值的编号
