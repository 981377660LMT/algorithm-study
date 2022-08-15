# https://www.acwing.com/problem/content/description/100/
# 类似分形之城

from typing import List

# 矩阵分治，按照mid分割成四块
# 希尔伯特曲线


def solve(n) -> List[List[int]]:
    if n == 1:
        return [[1]]
    if n == 2:
        return [[2, 3], [1, 4]]
    mid = n // 2
    pre = solve(mid)
    res = [[0] * n for _ in range(n)]
    counter = [0, n * n // 4, n * n // 2, 3 * n * n // 4]
    for i in range(mid):
        for j in range(mid):
            res[i][j] = pre[i][j] + counter[1]
            res[mid + i][j] = pre[mid - 1 - j][mid - 1 - i] + counter[0]
            res[i][mid + j] = pre[i][j] + counter[2]
            res[mid + i][mid + j] = pre[j][i] + counter[3]
    return res


n = int(input())
res = solve(n)
for i in range(n):
    print(res[i][0], end="")
    for j in range(1, n):
        print(" %d" % res[i][j], end="")
    print()
