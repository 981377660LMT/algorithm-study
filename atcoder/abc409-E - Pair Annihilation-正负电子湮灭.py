# abc409-E - Pair Annihilation-正负电子湮灭
#
# 给定一棵树。
# 每个顶点被赋予一个整数x，其含义如下：
# 若x>0，表示顶点上有x个正电子；
# 若x<0，表示顶点上有-x个负电子；
# 若x=0，表示顶点无粒子。
# 保证：所有顶点粒子数之和为0。
# 粒子移动规则沿边j移动一个正电子或电子需消耗能量wj；
# 当正电子与电子处于同一顶点时，会以相等数量成对湮灭。
# 目标计算湮灭所有粒子所需的最小总能量。

import sys
from typing import Tuple

sys.setrecursionlimit(int(1e6))

N = int(input())
X = list(map(int, input().split()))
E = [list(map(int, input().split())) for _ in range(N - 1)]

##

tree = [[] for _ in range(N)]
for u, v, w in E:
    u -= 1
    v -= 1
    tree[u].append((v, w))
    tree[v].append((u, w))


res = 0


def dfs(cur: int, pre: int) -> Tuple[int, int]:
    """树形dp, 向上返回负电子和正电子的数量."""
    global res
    a, b = 0, 0
    if X[cur] < 0:
        a += -X[cur]
    elif X[cur] > 0:
        b += X[cur]
    for next_, weight in tree[cur]:
        if next_ == pre:
            continue
        na, nb = dfs(next_, cur)
        res += weight * (na + nb)
        a += na
        b += nb
    min_ = min(a, b)
    a -= min_
    b -= min_
    return a, b


dfs(0, -1)

print(res)
