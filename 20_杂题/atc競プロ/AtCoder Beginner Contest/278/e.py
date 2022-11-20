from collections import defaultdict
import sys
from typing import List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 縦 H 行、横 W 列のグリッドがあります。 上から i 行目、左から j 行目のマスを (i,j) で表します。 (i,j) (1≤i≤H,1≤j≤W) には 1 以上 N 以下の整数 A
# i,j
# ​
#   が書かれています。

# 整数 h,w が与えられます。0≤k≤H−h,0≤l≤W−w を満たすすべての (k,l) の組について、次の問題を解いてください。

# k<i≤k+h,l<j≤l+w を満たす (i,j) を塗りつぶしたとき、塗りつぶされていないマスに書かれている数が何種類あるか求めよ。
# ただし、問題を解く際に実際にマスを塗りつぶすことはない（各問題が独立である）ことに注意してください。
# !没有被覆盖到的格子里的颜色数
if __name__ == "__main__":
    ROW, COL, n, windowRow, windowCol = map(int, input().split())
    grid = []
    for i in range(ROW):
        grid.append(list(map(int, input().split())))

    allCounter = defaultdict(int)
    for i in range(ROW):
        for j in range(COL):
            allCounter[grid[i][j]] += 1
    allCount = len(allCounter)

    # !二维滑动窗口
    def cal(row: int) -> List[int]:
        res, curFull = [], 0
        curCounter = defaultdict(int)
        for right in range(COL):
            for r in range(row, windowRow + row):
                curCounter[grid[r][right]] += 1
                if curCounter[grid[r][right]] == allCounter[grid[r][right]]:
                    curFull += 1
            if right >= windowCol:
                for r in range(row, windowRow + row):
                    curCounter[grid[r][right - windowCol]] -= 1
                    if (
                        curCounter[grid[r][right - windowCol]]
                        == allCounter[grid[r][right - windowCol]] - 1
                    ):
                        curFull -= 1
            if right >= windowCol - 1:
                res.append(allCount - curFull)
        return res

    for r in range(ROW - windowRow + 1):
        res = cal(r)
        print(*res)
