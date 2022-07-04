from itertools import product
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


DIR8 = [(1, 0), (0, 1), (-1, 0), (0, -1), (1, 1), (1, -1), (-1, 1), (-1, -1)]


def main() -> None:
    n = int(input())
    matrix = []
    for _ in range(n):
        row = [int(char) for char in input()]
        matrix.append(row)

    res = 0
    for sr, sc, di in product(range(n), range(n), DIR8):  # 八方旅人
        sb = []
        r, c = sr, sc
        for _ in range(n):
            sb.append(matrix[r][c])
            r, c = (r + di[0]) % n, (c + di[1]) % n
        cand = int("".join(map(str, sb)))
        res = max(res, cand)

    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
