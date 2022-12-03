from heapq import heappop, heappush
import sys
from typing import List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# (r,c) から (r,c+1) に移動する。A
# r,c
# ​
#   のコストがかかる。この移動は c<C のとき使える。
# (r,c) から (r,c−1) に移動する。A
# r,c−1
# ​
#   のコストがかかる。この移動は c>1 のとき使える。
# (r,c) から (r+1,c) に移動する。B
# r,c
# ​
#   のコストがかかる。この移動は r<R のとき使える。
# 1≤i<r を満たす整数 i を 1 つ選び、(r,c) から (r−i,c) に移動する。1+i のコストがかかる。


def sneaking(ROW: int, COL: int, cost1: List[List[int]], cost2: List[List[int]]) -> int:
    dist = [[INF] * COL for _ in range(ROW)]
    dist[0][0] = 0
    pq = [(0, 0, 0)]  # (cost,row,col)
    while pq:
        curDist, curRow, curCol = heappop(pq)
        if curRow == ROW - 1 and curCol == COL - 1:
            return curDist
        if dist[curRow][curCol] < curDist:
            continue
        if curCol + 1 < COL:
            cand = curDist + cost1[curRow][curCol]
            if dist[curRow][curCol + 1] > cand:
                dist[curRow][curCol + 1] = cand
                heappush(pq, (cand, curRow, curCol + 1))
        if curCol - 1 >= 0:
            cand = curDist + cost1[curRow][curCol - 1]
            if dist[curRow][curCol - 1] > cand:
                dist[curRow][curCol - 1] = cand
                heappush(pq, (cand, curRow, curCol - 1))
        if curRow + 1 < ROW:
            cand = curDist + cost2[curRow][curCol]
            if dist[curRow + 1][curCol] > cand:
                dist[curRow + 1][curCol] = cand
                heappush(pq, (cand, curRow + 1, curCol))

        for i in range(1, curRow + 1):
            cand = curDist + 1 + i
            if dist[curRow - i][curCol] > cand:
                dist[curRow - i][curCol] = cand
                heappush(pq, (cand, curRow - i, curCol))

    raise Exception("No path")


if __name__ == "__main__":
    ROW, COL = map(int, input().split())
    grid = [list(map(int, input().split())) for _ in range(ROW)]
    costs = [list(map(int, input().split())) for _ in range(ROW - 1)]
    print(sneaking(ROW, COL, grid, costs))
