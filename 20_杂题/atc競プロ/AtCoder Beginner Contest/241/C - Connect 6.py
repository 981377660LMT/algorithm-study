# !最多可以把两个格子变为白色 问能否形成六子棋(连续六个黑色)
# n<=1000

# !每个点枚举八个方向即可
# !起点、方向 然后while循环模拟

import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
DIR8 = [(1, 0), (1, -1), (0, -1), (-1, -1), (-1, 0), (-1, 1), (0, 1), (1, 1)]


def main() -> None:
    n = int(input())
    grid = []
    for _ in range(n):
        grid.append(list(input()))

    for r in range(n):
        for c in range(n):
            for dr, dc in DIR8:
                cr, cc = r, c
                count, remain = 1, 2 - int(grid[cr][cc] == ".")
                while True:
                    nr, nc = cr + dr, cc + dc
                    if nr < 0 or nr >= n or nc < 0 or nc >= n:
                        break
                    cr, cc = nr, nc
                    if grid[cr][cc] == ".":
                        if remain == 0:
                            break
                        remain -= 1
                    count += 1
                    if count == 6:
                        print("Yes")
                        exit(0)
    print("No")


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
