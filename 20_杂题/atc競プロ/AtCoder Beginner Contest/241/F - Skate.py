# 有一个H× W的迷宫(滑雪场)。网格上有N个障碍物，第i个的位置是(Xi，Yi)。
# 我们从(sx, sy)开始，每一步向上、下、左、右中的一个方向行走，
# 直到撞上障碍物，停在它前面的方格中。
# 求到达(gx,gy)所用的最少`步数`。若无法到达终点，输出-1。

# !ROW,COL<=1e9
# !N<=1e5
# !注意求步数而不是距离 因此是无权图求最短路 注意每次移动需要二分查找 不能线性移动

# !bfs忘记加visited了

from bisect import bisect_left, bisect_right
from collections import defaultdict, deque
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    ROW, COL, N = map(int, input().split())
    startRow, startCol = [int(num) - 1 for num in input().split()]
    endRow, endCol = [int(num) - 1 for num in input().split()]
    rowAdjMap = defaultdict(list)
    colAdjMap = defaultdict(list)
    for _ in range(N):
        r, c = map(int, input().split())
        r, c = r - 1, c - 1
        rowAdjMap[r].append(c)
        colAdjMap[c].append(r)

    for v in rowAdjMap.values():
        v.sort()
    for v in colAdjMap.values():
        v.sort()

    queue = deque([(0, startRow, startCol)])
    visited = set([(startRow, startCol)])
    while queue:
        dist, curRow, curCol = queue.popleft()
        if curRow == endRow and curCol == endCol:
            print(dist)
            exit(0)

        # !四个方向动 二分查找

        # 上
        up = bisect_right(colAdjMap[curCol], curRow) - 1
        if up >= 0:
            nr, nc = colAdjMap[curCol][up] + 1, curCol
            if 0 <= nr < ROW and 0 <= nc < COL and (nr, nc) not in visited:
                visited.add((nr, nc))
                queue.append((dist + 1, nr, nc))

        # 下
        down = bisect_left(colAdjMap[curCol], curRow)
        if down <= len(colAdjMap[curCol]) - 1:
            nr, nc = colAdjMap[curCol][down] - 1, curCol
            if 0 <= nr < ROW and 0 <= nc < COL and (nr, nc) not in visited:
                visited.add((nr, nc))
                queue.append((dist + 1, nr, nc))

        # 左
        left = bisect_right(rowAdjMap[curRow], curCol) - 1
        if left >= 0:
            nr, nc = curRow, rowAdjMap[curRow][left] + 1
            if 0 <= nr < ROW and 0 <= nc < COL and (nr, nc) not in visited:
                visited.add((nr, nc))
                queue.append((dist + 1, nr, nc))

        # 右
        right = bisect_left(rowAdjMap[curRow], curCol)
        if right <= len(rowAdjMap[curRow]) - 1:
            nr, nc = curRow, rowAdjMap[curRow][right] - 1
            if 0 <= nr < ROW and 0 <= nc < COL and (nr, nc) not in visited:
                visited.add((nr, nc))
                queue.append((dist + 1, nr, nc))

    print(-1)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
