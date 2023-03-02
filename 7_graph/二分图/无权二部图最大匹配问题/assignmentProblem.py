# 任务分配/任务调度
# 分配问题(assignment problem)

# n个人分配n项任务，一个人只能分配一项任务，一项任务只能分配给一个人，
# 将一项任务分配给一个人是需要支付报酬，如何分配任务，保证支付的报酬总数最小。
# 简单的说：就是n*n矩阵中，选取n个元素，每行每列各有一个元素，使得和最小。
# 1<=n<=500 -1e9<=Aij<=1e9
# https://judge.yosupo.jp/problem/assignment

from typing import List, Tuple

INF = int(1e18)


def hungarian(A: List[List[int]]) -> Tuple[int, List[int]]:
    """任务分配问题

    Args:
        A (List[List[int]]): A[i][j] 表示第i个人分配第j项任务的报酬
    Returns:
        Tuple[int, List[int]]: 最小报酬和,每一行选取第几列的元素
    """

    n = len(A) + 1
    m = len(A[0]) + 1
    P = [0] * m
    way = [0] * m
    U = [0] * n
    V = [0] * n
    for i in range(1, n):
        P[0] = i
        minV = [INF] * m
        used = [False] * m
        j0 = 0
        while P[j0] != 0:
            i0 = P[j0]
            j1 = 0
            used[j0] = True
            delta = INF
            for j in range(1, m):
                if used[j]:
                    continue
                if i0 == 0 or j == 0:
                    cur = -U[i0] - V[j]
                else:
                    cur = A[i0 - 1][j - 1] - U[i0] - V[j]
                if cur < minV[j]:
                    minV[j] = cur
                    way[j] = j0
                if minV[j] < delta:
                    delta = minV[j]
                    j1 = j
            for j in range(m):
                if used[j]:
                    U[P[j]] += delta
                    V[j] -= delta
                else:
                    minV[j] -= delta
            j0 = j1
        P[j0] = P[way[j0]]
        j0 = way[j0]
        while j0 != 0:
            P[j0] = P[way[j0]]
            j0 = way[j0]

    res = [-1] * (n - 1)
    for i in range(1, m):
        if P[i] != 0:
            res[P[i] - 1] = i - 1
    return -V[0], res


n = int(input())
A = [list(map(int, input().split())) for _ in range(n)]
cost, perm = hungarian(A)
print(cost)
print(*perm)


# 4 3 5
# 3 5 9
# 4 1 4
# 输出
# 9
# 2 0 1
