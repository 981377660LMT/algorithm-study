# https://atcoder.jp/contests/abc274/editorial/5098

# 监视器覆盖问题
# 给定一张 ROW行COL列的网格图。
# "."表示空地，"#"表示障碍物.
# 你可以在空白的位置放监控摄像头。监控有四个方向：前后左右。
# 一个方向的监控只能看到自己正方向的位置。
# 例如，向前的监控只能看到自己正前方的位置。
# 监控的视野会被障碍挡住。一个位置可以放多个监控。
# !你需要求出最少的放置数量，使得所有的空白位置都能被看到。
# 1 <= ROW, COL <= 300

# !每个空白由上下监控还是左右监控(行还是列覆盖它),监控要最少 => 最小点覆盖问题

from 匈牙利算法 import Hungarian

from typing import List


def securityCamera3(grid: List[str]) -> int:
    ROW, COL = len(grid), len(grid[0])
    H = Hungarian()
    for r in range(ROW):
        for c in range(COL):
            if grid[r][c] == "#":
                continue
            curRow, curCol = r, c
            while curRow > 0 and grid[curRow - 1][c] == ".":
                curRow -= 1
            while curCol > 0 and grid[r][curCol - 1] == ".":
                curCol -= 1
            # ! (cr,c) 和 (r,cc) 两个行/列监视的候选位置二选1
            cand1, cand2 = curRow * COL + c, r * COL + curCol
            H.addEdge(cand1, cand2)
    return len(H.work())


if __name__ == "__main__":
    ROW, COL = map(int, input().split())
    grid = [input() for _ in range(ROW)]
    print(securityCamera3(grid))
