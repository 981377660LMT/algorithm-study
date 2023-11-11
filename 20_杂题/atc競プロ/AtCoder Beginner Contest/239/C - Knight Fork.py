# 是否存在到两个点距离均为根号5的点

import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)

# 根号5
DIR8 = [(2, 1), (1, 2), (-1, 2), (-2, 1), (-2, -1), (-1, -2), (1, -2), (2, -1)]


def main() -> None:
    x1, y1, x2, y2 = map(int, input().split())
    res1 = [(x1 + dx, y1 + dy) for dx, dy in DIR8]
    res2 = [(x2 + dx, y2 + dy) for dx, dy in DIR8]
    print(["No", "Yes"][len(set(res1) & set(res2)) > 0])


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
