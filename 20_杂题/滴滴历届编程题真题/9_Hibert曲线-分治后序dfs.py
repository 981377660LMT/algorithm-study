from typing import List


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
ans = solve(n)
for i in range(n):
    print(ans[i][0], end='')
    for j in range(1, n):
        print(' %d' % ans[i][j], end='')
    print()

