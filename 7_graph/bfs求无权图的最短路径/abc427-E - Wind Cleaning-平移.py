# abc427-E - Wind Cleaning-平移(风来之国)
# https://atcoder.jp/contests/abc427/tasks/abc427_e
#
# 有一个 H 行 W 列的网格。从上数第 i 行、从左数第 j 列的格子记为 (i, j)。

# 高桥君站在这个网格的其中一个格子上，网格中的某些格子上则有垃圾。

# 每个格子的状态由 H 个长度为 W 的字符串 S_1, S_2, ..., S_H 给出。S_i 的第 j 个字符 S_{i,j} 代表了如下状态：

# S_{i,j} = 'T'：表示高桥君在格子 (i, j) 上。
# S_{i,j} = '#'：表示格子 (i, j) 上有垃圾。
# S_{i,j} = '.'：表示格子 (i, j) 是一个没有垃圾的空格子。
# 另外，高桥君所在的格子上没有垃圾。

# 高桥君可以重复进行以下操作：

# !选择上、下、左、右四个方向中的一个，将所有垃圾同时向该方向移动一格。如果垃圾移出网格，则垃圾消失。如果垃圾移动到高桥君所在的格子，高桥君就会被弄脏。
# 请判断高桥君是否能在不被弄脏的情况下清除所有垃圾。如果可以，请找出所需的最少操作次数。

# 约束条件

# 2 ≤ H, W ≤ 12
# S_i 是由 'T', '#', '.' 组成的长度为 W 的字符串。
# 恰好存在一个 (i, j) 使得 S_{i,j} = 'T'。
# 至少存在一个 (i, j) 使得 S_{i,j} = '#'。


from collections import deque
from typing import List


def solve(H: int, W: int, S: List[str]) -> int:
    t = -1
    flag = 0
    for i in range(H):
        for j in range(W):
            c = S[i][j]
            idx = i * W + j
            if c == "T":
                t = idx
            elif c == "#":
                flag |= 1 << idx

    if flag == 0:
        return 0

    total = (1 << (H * W)) - 1
    MASK = (1 << W) - 1
    row_base = [r * W for r in range(H)]

    def hit_t(mask):
        return (mask >> t) & 1

    def shift_up(m):
        nm = m >> W
        return -1 if hit_t(nm) else nm

    def shift_down(m):
        nm = (m << W) & total
        return -1 if hit_t(nm) else nm

    def shift_left(m):
        nm = 0
        for r in range(H):
            row = (m >> row_base[r]) & MASK
            row >>= 1
            nm |= row << row_base[r]
        return -1 if hit_t(nm) else nm

    def shift_right(m):
        nm = 0
        for r in range(H):
            row = (m >> row_base[r]) & MASK
            row = (row << 1) & MASK
            nm |= row << row_base[r]
        return -1 if hit_t(nm) else nm

    queue = deque([(flag, 0)])
    seen = {flag}
    while queue:
        m, d = queue.popleft()
        for f in (shift_up, shift_down, shift_left, shift_right):
            nm = f(m)
            if nm < 0:
                continue
            if nm == 0:
                return d + 1
            if nm not in seen:
                seen.add(nm)
                queue.append((nm, d + 1))
    return -1


if __name__ == "__main__":
    H, W = map(int, input().split())
    S = [input() for _ in range(H)]
    print(solve(H, W, S))
