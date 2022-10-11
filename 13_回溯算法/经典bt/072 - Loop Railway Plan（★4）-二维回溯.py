# 从源点出发四个方向乱走,每个格子最多走一次,直到走回源点
# !途中至少走三个格子
# 求经过的格子数量的最大值 如果不存在则返回-1

# 回溯
# ROW*COL<=16

# 黄金矿工
import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)


ROW, COL = map(int, input().split())
matrix = []
for _ in range(ROW):
    row = list(input())
    matrix.append(row)

res = -1
DIR4 = [[1, 0], [-1, 0], [0, 1], [0, -1]]


def bt(r: int, c: int, count: int, sr: int, sc: int) -> None:
    """注意要放在next里更新 因为每次进入都会被标记为访问过"""
    global res

    matrix[r][c] = "#"
    for dr, dc in DIR4:
        nr, nc = r + dr, c + dc
        if 0 <= nr < ROW and 0 <= nc < COL:
            if nr == sr and nc == sc and count + 1 >= 3:  # 回到原点
                res = max(res, count + 1)
            if matrix[nr][nc] == ".":
                bt(nr, nc, count + 1, sr, sc)
    matrix[r][c] = "."


for r in range(ROW):
    for c in range(COL):
        if matrix[r][c] == ".":
            bt(r, c, 0, r, c)

print(res)
