from bisect import bisect_right
from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# LRUD
# 2≤H,W≤1e9  二分查找
DIR4 = [(0, -1), (0, 1), (-1, 0), (1, 0)]
mp = {"L": (0, -1), "R": (0, 1), "U": (-1, 0), "D": (1, 0)}
if __name__ == "__main__":
    ROW, COL, sr, sc = map(int, input().split())
    sr, sc = sr - 1, sc - 1
    n = int(input())
    walls = []
    rowWall = defaultdict(list)
    colWall = defaultdict(list)
    for _ in range(n):
        x, y = map(int, input().split())
        x, y = x - 1, y - 1
        walls.append((x, y))
        rowWall[x].append(y)
        colWall[y].append(x)

    for v in rowWall.values():
        v.sort()
    for v in colWall.values():
        v.sort()

    # 预处理从每个格子出发,沿每个方向能走多远

    q = int(input())
    # i 番目の指示に対して高橋君は下記の行動を stepi回繰り返します。
    # 現在いるマスに対して、d iが表す向きに壁のないマスが隣接しているなら、
    # そのマスに移動する。 そのようなマスが存在しない場合は、何もしない。
    # print(upDp)
    for _ in range(q):
        dir, step = input().split()
        step = int(step)
        dr, dc = mp[dir]
        if dir == "R":
            # 右边第一个墙壁
            pos = bisect_right(rowWall[sr], sc)
            wallPos = rowWall[sr][pos] if pos < len(rowWall[sr]) else COL
            step = min(step, wallPos - 1 - sc)
            sr, sc = sr + dr * step, sc + dc * step
        elif dir == "L":
            # 左边第一个墙壁
            pos = bisect_right(rowWall[sr], sc)
            wallPos = rowWall[sr][pos - 1] if pos > 0 else -1
            step = min(step, sc - (wallPos + 1))
            sr, sc = sr + dr * step, sc + dc * step
        elif dir == "U":
            # 上边第一个墙壁
            pos = bisect_right(colWall[sc], sr)
            wallPos = colWall[sc][pos - 1] if pos > 0 else -1
            step = min(step, sr - (wallPos + 1))
            sr, sc = sr + dr * step, sc + dc * step
        elif dir == "D":
            # 下边第一个墙壁
            pos = bisect_right(colWall[sc], sr)
            wallPos = colWall[sc][pos] if pos < len(colWall[sc]) else ROW
            step = min(step, wallPos - 1 - sr)
            sr, sc = sr + dr * step, sc + dc * step
        # sr, sc = sr + dr * dist, sc + dc * dist
        # dist = min(step, dp[sr][sc])
        print(sr + 1, sc + 1)
