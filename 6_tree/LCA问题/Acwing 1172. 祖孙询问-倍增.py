# 给定一棵包含 n 个节点的有根无向树，节点编号互不相同，但不一定是 1∼n。
# 有 m 个询问，每个询问给出了一对节点的编号 x 和 y，询问 x 与 y 的祖孙关系。
# 对于每一个询问，若 x 是 y 的祖先则输出 1，若 y 是 x 的祖先则输出 2，否则输出 0。

# 2^15 为 32798
# 1≤n,m≤4×104,
# 1≤每个节点的编号≤4×104
from collections import defaultdict, deque

from typing import DefaultDict, List


def bfs(start: int) -> None:
    """dfs 3万左右会爆栈,Segmentation Fault，应该用bfs"""
    queue = deque([(start, -1, 0)])
    while queue:
        cur, parent, level = queue.popleft()
        levelMap[cur] = level
        parentMap[cur] = parent
        for next in adjMap[cur]:
            if next == parent:
                continue
            queue.append((next, cur, level + 1))


def getFa(parentMap: DefaultDict[int, int]) -> List[List[int]]:
    """nlogn预处理"""
    fa = [[0] * 16 for _ in range(MAX)]
    for i in range(MAX):
        fa[i][0] = parentMap[i]
    for j in range(15):
        for i in range(MAX):
            if fa[i][j] == -1:
                fa[i][j + 1] = -1
            else:
                fa[i][j + 1] = fa[fa[i][j]][j]
    return fa


def getLCA(root1: int, root2: int, level: DefaultDict[int, int], fa: List[List[int]]) -> int:
    """logn查询"""
    if level[root1] < level[root2]:
        root1, root2 = root2, root1

    # diff = level[root1] - level[root2]
    # bit = 0
    # while diff:
    #     if diff & 1:
    #         root1 = fa[root1][bit]
    #     bit += 1
    #     diff >>= 1

    # 二进制拼凑法
    # logn(下取整)
    for i in range(15, -1, -1):
        if level[fa[root1][i]] >= level[root2]:
            root1 = fa[root1][i]

    if root1 == root2:
        return root1

    for i in range(15, -1, -1):
        if fa[root1][i] != fa[root2][i]:
            root1 = fa[root1][i]
            root2 = fa[root2][i]

    # 再往上跳1步即可
    return fa[root1][0]


MAX = int(4e4 + 10)
n = int(input())
adjMap = defaultdict(set)
root = -1
for i in range(n):
    u, v = map(int, input().split())
    if v == -1:
        root = u
        continue
    adjMap[u].add(v)
    adjMap[v].add(u)


# 1. 利用bfs找到每个点的父亲和每个点的深度
levelMap, parentMap = defaultdict(lambda: -1), defaultdict(lambda: -1)
bfs(root)

# 2.处理出fa数组
fa = getFa(parentMap)

m = int(input())
for _ in range(m):
    root1, root2 = map(int, input().split())
    # 2. 上跳到相同高度,如果此时a,b相等，返回
    # 3. 如果不相等，那么同时向上跳，最后跳到LCA的子节点位置
    lca = getLCA(root1, root2, levelMap, fa)
    if lca == root1:
        print(1)
    elif lca == root2:
        print(2)
    else:
        print(0)

