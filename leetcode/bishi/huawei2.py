from collections import deque


ROLE = {
    0: [(0, 1), (0, -1), (-1, 0), (1, 0)],
    1: [(1, 2), (1, -2), (-1, 2), (-1, -2), (2, 1), (2, -1), (-2, 1), (-2, -1)],
}

ROW, COL = map(int, input().split())
grid = []  # '.':space  'X':obstacle 'S':changeRole
for _ in range(ROW):
    grid.append(list(input()))

sr, sc, er, ec = 0, 0, ROW - 1, COL - 1
queue = deque([(sr, sc, 0, 0)])  # (row, col, step, role)
visited = set([(sr, sc, 0)])  # (row, col, role)
while queue:
    curRow, curCol, step, role = queue.popleft()
    if curRow == er and curCol == ec:
        print(step)
        exit(0)

    if grid[curRow][curCol] == "S":
        nextRole = role ^ 1
        nextState = (curRow, curCol, nextRole)
        if nextState not in visited:
            queue.append((curRow, curCol, step + 1, nextRole))
            visited.add(nextState)

    for dr, dc in ROLE[role]:
        nr, nc = curRow + dr, curCol + dc
        if 0 <= nr < ROW and 0 <= nc < COL and grid[nr][nc] != "X":
            nextState = (nr, nc, role)
            if nextState not in visited:
                queue.append((nr, nc, step + 1, role))
                visited.add(nextState)

print(-1)
