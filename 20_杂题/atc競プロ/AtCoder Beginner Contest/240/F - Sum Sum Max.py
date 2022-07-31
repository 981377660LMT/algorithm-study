# 给定一些关系构造出C后，
# C的长度为m (m<=1e9) C由y1个x1 y2个x2 ... yn个xn组成
# B为C的前缀和，A为B的前缀和，问A的最大值

# !注意到A为关于i的二次函数 因此可以三分法求凸函数最值
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def cal():
    n, m = map(int, input().split())
    pair = [map(int, input().split()) for _ in range(n)]


def main() -> None:
    T = int(input())
    for _ in range(T):
        cal()


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
