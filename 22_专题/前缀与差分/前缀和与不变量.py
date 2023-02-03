# 对某种数组上的操作,
# 可以前缀和数组/差分数组(不同的域上)
# https://atcoder.jp/contests/arc119/tasks/arc119_c

# 输入 n(2≤n≤3e5) 和长为 n 的数组 a(1≤a[i]≤1e9)。

# !每次操作，你可以选择两个相邻的数字，把它们都加一，或者都减一。
# 对于 a 的一个连续子数组 b，如果可以通过执行任意多次操作，
# 使 b 的所有元素为 0，则称 b 为好子数组。
# 输出 a 的好子数组的数量。

# 一次会改变多个数的题目，往往入手点在「不变量」上，也就是操作不会改变什么
# !通过奇偶前缀和(交错和)构造不变量
# !那么题目就变成有多少个区间的交错和等于 0
# 这是个经典问题，做法与「和为 0 的子数组个数」是一样的，用前缀和 + 哈希表解决。

from collections import defaultdict
from typing import List


def arcWrecker2(nums: List[int]) -> int:
    preSum = defaultdict(int, {0: 1})  # 如果记录索引就是{0: -1}
    res, curSum = 0, 0
    for i, num in enumerate(nums):
        curSum += num if i & 1 else -num
        if curSum in preSum:
            res += preSum[curSum]
        preSum[curSum] += 1
    return res


if __name__ == "__main__":

    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n = int(input())
    nums = list(map(int, input().split()))
    print(arcWrecker2(nums))
