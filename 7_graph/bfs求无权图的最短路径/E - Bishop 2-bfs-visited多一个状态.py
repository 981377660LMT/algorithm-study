"""bfs的思想 visited需要多存一个方向的状态"""

# 给出一个n* n的地图，其中'.'是平地，其中'#'是障碍，并给出一个起点和终点，在不跨越障碍的情况下，
# 可以往左上，右上，右下，左下四个斜角方向一次移动任意个单位，
# 询问从起点移动到终点的最小步数，若不存在路径，则输出-1。
# 2≤N≤1500


# !修改正常BFS的拓展方式即可，每次拓展的时候把四个斜角方向可以拓展的点全部拓展。
# !注意dijk取等号时的情况不要break 还要继续看这个方向
# !注意每次移动可以移动任意个单位

# TLE
# !因为无权图 所以 queue 可以换 queue

from collections import deque
import sys
import os


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")


DIR4 = ((1, 1), (1, -1), (-1, -1), (-1, 1))


def main() -> None:
    n = int(input())
    ROW, COL = n, n
    sr, sc = map(int, input().split())
    sr, sc = sr - 1, sc - 1
    er, ec = map(int, input().split())
    er, ec = er - 1, ec - 1

    grid = []
    for _ in range(ROW):
        row = list(input())
        grid.append(row)

    visited = [[[False] * 4 for _ in range(COL)] for _ in range(ROW)]
    visited[0][0] = [True] * 4
    queue = deque([(0, sr, sc)])
    while queue:
        step, r, c = queue.popleft()
        if r == er and c == ec:
            print(step)
            exit(0)
        for di, (dr, dc) in enumerate(DIR4):
            nr, nc = r + dr, c + dc
            while True:
                if (
                    0 <= nr < ROW
                    and 0 <= nc < COL
                    and grid[nr][nc] == "."
                    and not visited[nr][nc][di]
                ):
                    visited[nr][nc][di] = True
                    queue.append((step + 1, nr, nc))
                else:
                    break
                nr, nc = nr + dr, nc + dc

    print(-1)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
