from collections import deque
from typing import Deque, List


optCount = int(input())

mat = []
strings: List[Deque] = []
# 格子颜色rgb
for _ in range(8):
    chars = list(input())
    mat.append(chars)

for _ in range(8):
    chars = list(input())
    strings.append(deque(chars))


opt = []
for _ in range(optCount):
    # 1<=x,y<=8
    x, y, direct = list(input())
    x, y = int(x) - 1, int(y) - 1
    opt.append((x, y, direct))

res = []


def main():
    def dfs(mat: List[List[str]], color: str, x: int, y: int) -> None:
        if mat[x][y] == color:
            mat[x][y] = '*'
        for dx, dy in [(1, 0), (-1, 0), (0, 1), (0, -1)]:
            nx, ny = x + dx, y + dy
            if 0 <= nx < 8 and 0 <= ny < 8 and mat[nx][ny] == color:
                dfs(mat, color, nx, ny)

    def move(mat: List[List[str]], direction: str) -> None:

        # wasd
        if direction == 's':
            for col in range(8):
                write = 7
                for read in range(7, -1, -1):
                    if mat[read][col] != '*':
                        mat[write][col] = mat[read][col]
                        write -= 1
                for row in range(write, -1, -1):
                    mat[row][col] = '*'
        elif direction == 'w':
            for col in range(8):
                write = 0
                for read in range(8):
                    if mat[read][col] != '*':
                        mat[write][col] = mat[read][col]
                        write += 1
                for row in range(write, 8):
                    mat[row][col] = '*'
        elif direction == 'a':
            for row in range(8):
                write = 0
                for read in range(8):
                    if mat[row][read] != '*':
                        mat[row][write] = mat[row][read]
                        write += 1
                for col in range(write, 8):
                    mat[row][col] = '*'
        elif direction == 'd':
            for row in range(8):
                write = 7
                for read in range(7, -1, -1):
                    if mat[row][read] != '*':
                        mat[row][write] = mat[row][read]
                        write -= 1
                for col in range(write, -1, -1):
                    mat[row][col] = '*'

    def fill(mat: List[List[str]], direction: str) -> int:
        res = 0
        if direction in ('a', 'd'):
            # 依次填补列
            for col in range(8):
                curString = strings[col]
                for row in range(8):
                    if mat[row][col] == '*':
                        res += 1
                        mat[row][col] = '*' if curString else curString.popleft()
        else:
            # 依次填补列
            for row in range(8):
                curString = strings[row]
                for col in range(8):
                    if mat[row][col] == '*':
                        res += 1
                        mat[row][col] = '*' if curString else curString.popleft()
        return res

    for i in range(optCount):
        x, y, direct = opt[i]
        # 消除
        dfs(mat, mat[x][y], x, y)
        # 移动
        move(mat, direct)
        # 补全
        count = fill(mat, direct)
        res.append(count)


main()

for num in res:
    print(num)

