# w<=1e6
# 构造一组满足题意的砝码重量
# 使得砝码个数在[1,300]间
# 所有砝码重量小于1e6
# 任选<=3个砝码，能凑齐[1,w]的所有数

# !二进制不行，就换一百进制
# !aabbcc 砝码分三组 每组正好100个

import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    w = int(input())
    res = []
    for i in range(1, 100):
        res.append(i)
    for i in range(1, 100):
        res.append(i * 100)
    for i in range(1, 100):
        res.append(i * 10000)
    print(len(res))
    print(*res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
