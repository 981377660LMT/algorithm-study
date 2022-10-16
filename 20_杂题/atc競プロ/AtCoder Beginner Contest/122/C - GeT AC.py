# 求区间内子串出现的次数

# !先找到所有的起始点 然后前缀和区间计数
from itertools import accumulate
import sys
from typing import List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def findAll(string: str, target: str) -> List[int]:
    """找到所有匹配的字符串起始位置"""
    start = 0
    res = []
    while True:
        pos = string.find(target, start)
        if pos == -1:
            break
        else:
            res.append(pos)
            start = pos + 1

    return res


if __name__ == "__main__":

    n, q = map(int, input().split())
    s = input()
    indexes = findAll(s, "AC")

    # !前缀和求区间内子串出现个数
    nums = [0] * n
    for index in indexes:
        nums[index] = 1
    preSum = [0] + list(accumulate(nums))
    for _ in range(q):
        left, right = map(int, input().split())
        left, right = left - 1, right - 1
        print(preSum[right] - preSum[left])
