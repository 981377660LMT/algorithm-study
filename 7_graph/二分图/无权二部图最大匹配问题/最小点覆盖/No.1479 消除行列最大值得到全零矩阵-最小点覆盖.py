# 消除行列最大值得到全零矩阵
# !每次操作可以消除一行或一列中的所有最大值，
# 问最少需要多少次操作可以消除所有元素，使得矩阵中的所有元素都为0。
# ROW,COL<=500 0<=Aij<=5e5

# 最小点覆盖问题
# !每一个不为0的数,`要么由行消除,要么由列消除`
# 选出最少的行+列,使得每个不为0的数都被消除 => 最小点覆盖问题

from 匈牙利算法 import Hungarian

from collections import defaultdict
from typing import List


def solve(grid: List[List[int]]) -> int:
    ROW, COL = len(grid), len(grid[0])
    mp = defaultdict(list)
    for i in range(ROW):
        for j in range(COL):
            mp[grid[i][j]].append((i, j))

    res = 0
    for v, edges in mp.items():
        if v == 0:
            continue
        H = Hungarian()
        id1, id2 = dict(), dict()
        for u, v in edges:
            id1.setdefault(u, len(id1))
            id2.setdefault(v, len(id2))
            H.addEdge(id1[u], id2[v])
        res += len(H.work())

    return res


if __name__ == "__main__":
    n, m = map(int, input().split())
    grid = [list(map(int, input().split())) for _ in range(n)]
    print(solve(grid))
