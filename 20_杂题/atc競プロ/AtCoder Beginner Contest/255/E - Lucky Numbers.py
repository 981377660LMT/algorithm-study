"""
n<=1e5 m<=10 暗示时间复杂度O(nm)
数组nums两两相邻的数字之和为pairSums[i]
这个数组最多有几个lucky数字

pairSums的性质
奇数切片一起偏移 偶数切片一起偏移
计算偏移量
"""
from collections import defaultdict
import sys
import os

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)


def main() -> None:
    n, m = map(int, input().split())
    pairSums = list(map(int, input().split()))
    luckies = list(map(int, input().split()))
    evenOffset, oddOffset = defaultdict(int), defaultdict(int)
    evenOffset[0] = 1  # 首项为基准0，计算偏移量 注意奇偶组偏移方向相反

    pre = 0
    for i in range(n - 1):
        pre = pairSums[i] - pre
        if i & 1:
            oddOffset[pre] += 1
        else:
            evenOffset[pre] += 1

    res = defaultdict(int)  # 处理组内数字与每个幸运数的偏移量
    for lucky in luckies:
        for num, count in evenOffset.items():
            res[num - lucky] += count
        for num, count in oddOffset.items():
            res[lucky - num] += count

    print(max(res.values(), default=0))


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            try:
                main()
            except (EOFError, ValueError):
                break
    else:
        main()
