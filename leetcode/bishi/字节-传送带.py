# U，D，L，R 分别表示朝向[上，下，左，右] 的传送带，
# 站在传送带上的人会被强制移动到其指向的下一个位置
# 如果下一个位置还是传送带，会被继续传下去
# 如果传送带指向迷宫外，玩家会撞在墙上，昏过去，游戏结束，无法再到达出口
# O表示迷宫出口
# !请你找到有多少点按照规则行走不能到达终点。


# !1.反向思考，把传送带视为一个单向板
# 从终点反向遍历，假设你处于一个空地上，
# !你左侧的点只有为R的时候才可以从你这个位置到达R那个位置，然后查找地图上所有O的数量。

# !2.注意:传送带不要一次性while模拟走完 要bfs一个一个看

from collections import deque

ROW, COL = map(int, input().split())
grid = [list(input()) for _ in range(ROW)]
er, ec = -1, -1
for r in range(ROW):
    for c in range(COL):
        if grid[r][c] == "O":
            er, ec = r, c
            break
    if er != -1:
        break

queue = deque([(er, ec)])
while queue:
    curRow, curCol = queue.popleft()
    # 上
    if curRow - 1 >= 0 and grid[curRow - 1][curCol] in "D.":
        queue.append((curRow - 1, curCol))
        grid[curRow - 1][curCol] = "O"

    # 下
    if curRow + 1 < ROW and grid[curRow + 1][curCol] in "U.":
        queue.append((curRow + 1, curCol))
        grid[curRow + 1][curCol] = "O"

    # 左
    if curCol - 1 >= 0 and grid[curRow][curCol - 1] in "R.":
        queue.append((curRow, curCol - 1))
        grid[curRow][curCol - 1] = "O"

    # 右
    if curCol + 1 < COL and grid[curRow][curCol + 1] in "L.":
        queue.append((curRow, curCol + 1))
        grid[curRow][curCol + 1] = "O"


good = sum(row.count("O") for row in grid[er])
print(ROW * COL - good)
