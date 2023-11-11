"""
row,col<=10
判断 B 是不是 A 的子矩阵(选择某些行列构成的子矩阵)

combinations选择行列计数
时间复杂度 O(C(10,5)^2*100)
"""


from itertools import combinations
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


ROW1, COL1 = map(int, input().split())
grid1 = []
for _ in range(ROW1):
    row = list(map(int, input().split()))
    grid1.append(row)

ROW2, COL2 = map(int, input().split())
grid2 = []
for _ in range(ROW2):
    row = list(map(int, input().split()))
    grid2.append(row)

target = [num for row in grid2 for num in row]

for rows in combinations(range(ROW1), ROW2):
    for cols in combinations(range(COL1), COL2):
        if list(grid1[r][c] for r in rows for c in cols) == target:
            print("Yes")
            exit(0)
print("No")

# 不太好的O(4^n*n^2)
# for state1 in range(1 << ROW1):
#     for state2 in range(1 << COL1):
#         cur = []
#         for r in range(ROW1):
#             for c in range(COL1):
#                 if (state1 >> r) & 1 and (state2 >> c) & 1:
#                     cur.append(grid1[r][c])

#         if len(cur) == len(target) and cur == target:
#             print("Yes")
#             exit()

# print("No")
