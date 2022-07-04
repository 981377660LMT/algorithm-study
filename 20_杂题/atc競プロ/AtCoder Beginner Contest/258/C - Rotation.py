# 两种查询
# 1.将字符串向右旋转x次
# 2.输出第x个字符
# n<=5e5 Q<=5e5 x<=5e5

# !显然不能模拟 而是记录偏移量


import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n, q = map(int, input().split())
    s = input()
    offset = 0
    for _ in range(q):
        qType, moveOrIndex = map(int, input().split())
        if qType == 1:
            move = moveOrIndex
            offset = (offset + move) % n
        else:
            index = moveOrIndex - 1
            print(s[(index - offset) % n])


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
