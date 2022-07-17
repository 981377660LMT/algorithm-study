# 求矩形第四个顶点的坐标 矩形的边平行于坐标轴
# !找到只出现一次的横/纵 坐标即可 异或

import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    points = [tuple(map(int, input().split())) for _ in range(3)]
    px = [x for x, _ in points]
    py = [y for _, y in points]

    # !异或找到只出现一次的横/纵 坐标
    print(px[0] ^ px[1] ^ px[2], py[0] ^ py[1] ^ py[2])


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
