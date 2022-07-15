# 倒水问题的加强版

# 给定三个整数a,b,c，
# 现在有两种操作，
# 第一种是选择两个数同时-1，
# 第二种是三个数同时-1，
# 问你最少多少次可以把他们变成全0

# a,b,c<=1e18


# 假设a>b>c， 然后如果a>b+c那么肯定是不行的，
# 否则我们肯定可以通过第二种操作使得他们变成a = b+c的形式


import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    a, b, c = sorted(map(int, input().split()))
    if c > a + b:
        print(-1)
        exit(0)

    diff = a + b - c
    if diff > a:
        print(-1)
        exit(0)

    print(c)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
