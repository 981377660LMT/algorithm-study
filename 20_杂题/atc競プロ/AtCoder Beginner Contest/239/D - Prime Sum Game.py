# 小T先在 A∼B 中选择一个数，小A再在 C∼D 中选择一个数，如果两个数之和为质数则小A赢，否则小T赢。
# 1<=A,B,C,D<=100
# !直接枚举小T的每种取值能不能赢即可

import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)


isPrime = [True] * 201
for i in range(2, 201):
    if not isPrime[i]:
        continue
    for j in range(i * i, 201, i):
        isPrime[j] = False


def main() -> None:
    a, b, c, d = map(int, input().split())
    for num1 in range(a, b + 1):
        if not any(isPrime[num1 + num2] for num2 in range(c, d + 1)):
            print("Takahashi")
            exit(0)
    print("Aoki")


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
