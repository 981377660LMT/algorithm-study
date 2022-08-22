# from collections import deque
# import sys
# from typing import Counter

# sys.setrecursionlimit(int(1e9))
# input = lambda: sys.stdin.readline().rstrip("\r\n")
# MOD = 998244353
# INF = int(4e18)

# n = int(input())

# bad = set(["q", "w", "k", "h", "j"])

# # owiq jjww 不是
# def check1(s: str) -> bool:
#     count = 0
#     for char in s:
#         if char in bad:
#             count += 1
#             if count >= 5:
#                 return True
#         else:
#             count = 0
#     return False


# queue = deque([])
# counter = Counter()


# def check2(s: str) -> bool:
#     queue.append(s)
#     res = False
#     if len(queue) >= 10:
#         counter[queue.popleft()] -= 1
#     if counter[s]:
#         res = True
#     counter[s] += 1
#     return res


# for _ in range(n):
#     s = input()
#     if check1(s) or check2(s):
#         print("yes")
#     else:
#         print("no")
# # no
# # no
# # no
# # yes
# # yes
# # no
# # no
# # no
# # no
# # yes
# # no
# # no
# # no

# # no
# # no
# # no
# # yes
# # yes
# # no
# # no
# # no
# # yes
# # yes
# # no
# # no
# # no
##############################################################

DIR4 = ((-1, 0), (1, 0), (0, 1), (0, -1))

# 反着走
MAPPING = {"U": (-1, 0), "D": (1, 0), "L": (0, -1), "R": (0, 1)}
# 迷宫中多少个位置 玩家不能到达出口

from collections import defaultdict, deque
import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")


ROW, COL = map(int, input().split())
grid = []
startRow, startCol = -1, -1
for _ in range(ROW):
    row = list(input())
    grid.append(row)


for r in range(ROW):
    flag = False
    for c in range(COL):
        if grid[r][c] == "O":
            startRow, startCol = r, c
            flag = True
            break
    if flag:
        break


# 1. 成环 检测哪些传送带在环上(组成环)
trans = set()  # 所有的传送带
adjMap = defaultdict(list)
deg = defaultdict(int)
for r in range(ROW):
    for c in range(COL):
        cur = r * COL + c
        if grid[r][c] in MAPPING:
            trans.add(cur)
            dr, dc = MAPPING[grid[r][c]]
            nr, nc = r + dr, c + dc
            if 0 <= nr < ROW and 0 <= nc < COL and grid[nr][nc] in MAPPING:
                next = nr * COL + nc
                trans.add(next)
                adjMap[cur].append(next)
                deg[next] += 1

queue = deque([p for p in trans if deg[p] == 0])

while queue:
    cur = queue.popleft()
    for next in adjMap[cur]:
        deg[next] -= 1
        if deg[next] == 0:
            queue.append(next)
onCycle = set([(p // COL, p % COL) for p in trans if deg[p] != 0])

# !从起点能走到多少个有效的格子
ok = set([(startRow, startCol)])
visited = set([(startRow, startCol)]) | onCycle
queue = deque([(startRow, startCol)])

while queue:
    curRow, curCol = queue.popleft()
    if not (0 <= curRow < ROW and 0 <= curCol < COL):
        continue

    # 正常格子
    if grid[curRow][curCol] not in MAPPING:
        ok.add((curRow, curCol))
        for dr, dc in DIR4:
            nr, nc = curRow + dr, curCol + dc
            if (0 <= nr < ROW and 0 <= nc < COL) and (nr, nc) not in visited:
                visited.add((nr, nc))
                queue.append((nr, nc))

    # !传送带能否传送至正常格子
    else:
        sRow, sCol = curRow, curCol
        isOk = False
        curVisted = set([(sRow, sCol)])
        while True:
            if grid[sRow][sCol] not in MAPPING:
                isOk = True
                break
            else:
                dr, dc = MAPPING[grid[curRow][curCol]]
                nr, nc = sRow + dr, sCol + dc
                if (0 <= nr < ROW and 0 <= nc < COL) and (nr, nc) not in curVisted:
                    sRow, sCol = nr, nc
                    curVisted.add((nr, nc))
                else:
                    isOk = False
                    break

        # visited |= curVisted
        if isOk:
            ok |= curVisted


print(ROW * COL - len(ok))
