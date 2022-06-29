# 求最大的子(序列)矩阵 矩阵每个数相同
# ROW<=8 COL<=10000
# !二进制枚举可能的行即可 这些行在某列上都是相同的

from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)


ROW, COL = map(int, input().split())
matrix = []
for _ in range(ROW):
    matrix.append(list(map(int, input().split())))

res = 0
for state in range(1, 1 << ROW):
    rows, rowCount = [], 0
    for r in range(ROW):
        if state & (1 << r):
            rows.append(r)
            rowCount += 1

    counter = defaultdict(int)  # 每个值有几列相同
    for c in range(COL):
        if len(set(matrix[r][c] for r in rows)) == 1:
            counter[matrix[rows[0]][c]] += 1

    cand = max(counter.values(), default=0) * rowCount  # !max 必加 default
    if cand > res:
        res = cand


print(res)
