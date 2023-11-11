"""
n<=1e5 m<=10 暗示时间复杂度O(nm)
数组nums两两相邻的数字之和为pairSums[i]
这个数组最多有几个lucky数字
"""

# https://atcoder.jp/contests/abc255/tasks/abc255_e
# https://www.luogu.com.cn/blog/endlesscheng/solution-at-abc255-e

# 输入 n (2≤n≤1e5) 和 m(≤10)，长为 n-1 的数组 pairSums 和长为 m 的严格递增数组 luckies，
# 元素值范围在 [-1e9,1e9]。
# 数组 luckies 中的元素叫做幸运数。
# !对于一个长为 n 的序列 A，如果所有相邻元素之和满足 A[i]+A[i+1]=pairSums[i]，则称为一个好序列。
# 输出好序列中最多能有多少个数是幸运数（重复数字也算，见样例）。

# !确定序列 A 中的任意一个数，就能确定整个序列 A。
# !以首项为参照物，计算奇数组和偶数组的偏移量
from typing import List


def luckyNumbers(n: int, pairSums: List[int], luckies: List[int]) -> int:
    counter0, counter1 = defaultdict(int), defaultdict(int)
    counter0[0] = 1  # !首项为基准0，计算偏移量 注意奇偶组偏移方向相反

    pre = 0
    for i in range(n - 1):
        pre = pairSums[i] - pre
        if i & 1:
            counter1[pre] += 1
        else:
            counter0[pre] += 1

    offset = defaultdict(int)  # 处理组内数字与每个幸运数的偏移量
    for lucky in luckies:
        for num, count in counter0.items():
            offset[num - lucky] += count
        for num, count in counter1.items():
            offset[lucky - num] += count

    return max(offset.values(), default=0)


if __name__ == "__main__":
    from collections import defaultdict
    import sys

    sys.setrecursionlimit(int(1e6))
    input = sys.stdin.readline
    n, _ = map(int, input().split())
    pairSums = list(map(int, input().split()))
    luckies = list(map(int, input().split()))
    print(luckyNumbers(n, pairSums, luckies))
