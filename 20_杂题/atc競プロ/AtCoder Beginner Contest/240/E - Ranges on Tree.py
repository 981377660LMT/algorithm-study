# dfs序区间的性质
# 如果子树有子集关系 那么对应区间也有子集关系
# 如果子树没有交集 那么对应区间也没有交集

# !dfs序映射到区间上 然后每个根找到他最左边的叶子结点和最右边的叶子结点


import sys
import os
from typing import Tuple

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)

# !dfs超时了


def main1() -> None:
    def dfs(cur: int, pre: int) -> Tuple[int, int]:
        """返回子树中 (最左边的叶子结点, 最右边的叶子结点)"""
        nonlocal dfsId

        if len(adjList[cur]) == 1 and adjList[cur][0] == pre:  # 叶子结点
            lefts[cur], rights[cur] = dfsId, dfsId
            dfsId += 1
            return lefts[cur], rights[cur]

        for next in adjList[cur]:
            if next == pre:
                continue
            nLeft, nRight = dfs(next, cur)
            lefts[cur], rights[cur] = min(lefts[cur], nLeft), max(rights[cur], nRight)

        return lefts[cur], rights[cur]

    n = int(input())
    adjList = [[] for _ in range(n)]
    for _ in range(n - 1):
        a, b = map(int, input().split())
        a, b = a - 1, b - 1
        adjList[a].append(b)
        adjList[b].append(a)

    lefts = [int(1e9)] * n  # 子树的最左边的叶子结点
    rights = [-int(1e9)] * n  # 子树的最右边的叶子结点
    dfsId = 1

    dfs(0, -1)
    for i in range(n):
        print(lefts[i], rights[i])


# !沿着dfs序倒着从下往上dp
def main() -> None:

    n = int(input())
    adjList = [[] for _ in range(n)]
    for _ in range(n - 1):
        a, b = map(int, input().split())
        a, b = a - 1, b - 1
        adjList[a].append(b)
        adjList[b].append(a)

    left = [int(1e9)] * n  # 子树的最左边的叶子结点
    right = [-int(1e9)] * n  # 子树的最右边的叶子结点
    dfsId = 1
    order = []  # dfs序
    parent = [-1] * n

    queue = [(0, -1)]
    while queue:
        cur, pre = queue.pop()  # dfs遍历
        order.append(cur)
        parent[cur] = pre
        if len(adjList[cur]) == 1 and adjList[cur][0] == pre:  # 叶子结点
            left[cur], right[cur] = dfsId, dfsId
            dfsId += 1
            continue

        for next in adjList[cur]:
            if next == pre:
                continue
            queue.append((next, cur))

    for root in order[::-1]:
        pre = parent[root]
        if pre != -1:
            left[pre] = min(left[pre], left[root])
            right[pre] = max(right[pre], right[root])

    for i in range(n):
        print(left[i], right[i])


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
