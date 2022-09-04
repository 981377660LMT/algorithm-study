# 华为机试-穿过障碍的最短时间-状压+最短路

# 0: 路
# 1: 墙
# 2: 起点
# 3: 终点
# 4: 陷阱
# 6: 炸弹

# 每走一步需要1秒
# 走到陷阱上需要3秒
# 走到炸弹上会激活炸弹将上下左右的墙炸成路 注意不能炸掉陷阱
# 状压+最短路
# 864. 获取所有钥匙的最短路径-多维度bfs

from collections import defaultdict
from heapq import heappop, heappush

INF = int(1e18)

DIR4 = ((0, 1), (1, 0), (-1, 0), (0, -1))
ROW, COL = map(int, input().split())
grid = []
for _ in range(ROW):
    row = list(map(int, input().split()))
    grid.append(row)

sr, sc = -1, -1
tr, tc = -1, -1
bomb_id = dict()
bomb_count = 0
trap_id = dict()
trap_count = 0
for r in range(ROW):
    for c in range(COL):
        if grid[r][c] == 2:
            sr, sc = r, c
        if grid[r][c] == 3:
            tr, tc = r, c
        if grid[r][c] == 4:
            trap_id[(r, c)] = trap_count
            trap_count += 1
        if grid[r][c] == 6:
            bomb_id[(r, c)] = bomb_count
            bomb_count += 1


def is_broken_wall(r: int, c: int, bomb_state: int) -> bool:
    bomb = 0
    for dr, dc in DIR4:
        nr, nc = r + dr, c + dc
        if 0 <= nr < ROW and 0 <= nc < COL and grid[nr][nc] == 6:
            bomb |= 1 << bomb_id[(nr, nc)]
    return bomb_state & bomb != 0


pq = [(0, sr, sc, 0, 0)]  # (time, r, c, bomb_state, trap_state)
dist = defaultdict(lambda: INF)  # (r, c, bomb_state, trap_state)
dist[(sr, sc, 0, 0)] = 0
while pq:
    cur_cost, cur_row, cur_col, cur_bomb_state, cur_trap_state = heappop(pq)
    if cur_row == tr and cur_col == tc:
        print(cur_cost)
        exit(0)

    for dr, dc in DIR4:
        nr, nc = cur_row + dr, cur_col + dc
        if not (0 <= nr < ROW and 0 <= nc < COL):
            continue

        if grid[nr][nc] in (0, 2, 3):
            cand = cur_cost + 1
            if cand < dist[(nr, nc, cur_bomb_state, cur_trap_state)]:
                dist[(nr, nc, cur_bomb_state, cur_trap_state)] = cand
                heappush(pq, (cand, nr, nc, cur_bomb_state, cur_trap_state))
        elif grid[nr][nc] == 1:  # 墙是否被摧毁
            if is_broken_wall(nr, nc, cur_bomb_state):
                cand = cur_cost + 1
                if cand < dist[(nr, nc, cur_bomb_state, cur_trap_state)]:
                    dist[(nr, nc, cur_bomb_state, cur_trap_state)] = cand
                    heappush(pq, (cand, nr, nc, cur_bomb_state, cur_trap_state))
        elif grid[nr][nc] == 4:
            cand = cur_cost + (1 if cur_trap_state & (1 << trap_id[(nr, nc)]) else 3)
            next_trap_state = cur_trap_state | (1 << trap_id[(nr, nc)])
            if cand < dist[(nr, nc, cur_bomb_state, next_trap_state)]:
                dist[(nr, nc, cur_bomb_state, next_trap_state)] = cand
                heappush(pq, (cand, nr, nc, cur_bomb_state, next_trap_state))
        elif grid[nr][nc] == 6:
            cand = cur_cost + 1
            next_bomb_state = cur_bomb_state | (1 << bomb_id[(nr, nc)])
            if cand < dist[(nr, nc, next_bomb_state, cur_trap_state)]:
                dist[(nr, nc, next_bomb_state, cur_trap_state)] = cand
                heappush(pq, (cand, nr, nc, next_bomb_state, cur_trap_state))
