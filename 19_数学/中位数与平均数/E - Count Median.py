# https://atcoder.jp/contests/abc169/tasks/abc169_e

# 输入 n (2≤n≤2e5) 和 n 个区间的左右端点，区间范围在 [1,1e9]。
# !每个区间内选一个整数，然后计算这 n 个整数的中位数。你能得到多少个不同的中位数？
# 注：偶数长度的中位数是中间两个数的平均值（没有下取整）。

from typing import List


def cal(sortedNums: List[int]) -> int:
    """如果长度为偶数，返回中间两个数的和；如果长度为奇数，返回中位数。"""
    n = len(sortedNums)
    if n & 1:
        return sortedNums[n // 2]
    return sortedNums[n // 2 - 1] + sortedNums[n // 2]


def countMedian(intervals: List[List[int]]) -> int:
    lowers, uppers = [], []
    for left, right in intervals:
        lowers.append(left)
        uppers.append(right)
    lowers.sort()
    uppers.sort()
    minMid = cal(lowers)
    maxMid = cal(uppers)
    return maxMid - minMid + 1


if __name__ == "__main__":
    n = int(input())
    intervals = list(list(map(int, input().split())) for _ in range(n))
    print(countMedian(intervals))
