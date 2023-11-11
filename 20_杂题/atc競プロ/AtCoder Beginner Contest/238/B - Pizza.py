# 切披萨，先在 12:00 的位置(钟表的位置) 切一刀，然后按照给定的序列 A ，
# 每次先顺时针旋转 nums[i]  度，然后在在 12:00 的位置切一刀，
# 问最后的所有披萨块中圆心角最大的是多少度？


import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def main() -> None:
    n = int(input())
    nums = list(map(int, input().split()))
    preSum = [0]
    for num in nums:
        preSum.append((preSum[-1] + num) % 360)
    preSum.append(360)
    splits = sorted(set(preSum))
    cands = [abs(pre - cur) for pre, cur in zip(splits, splits[1:])]
    print(max(cands))


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
