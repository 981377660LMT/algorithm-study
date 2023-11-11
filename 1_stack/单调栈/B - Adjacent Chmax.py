# 给定一个1-n的全排列 可以每次把相邻两个元素替换为最大值
# !求最后可能的数组个数
# n<=5000

# !原问题是dp (从之前替换的状态转移,状态为left,right)
# !对每个值求出作为最大值的范围

# !dp[index][right] 表示前index个数中 最大值右端点到right时的个数

# https://atcoder.jp/contests/agc058/submissions/34054055


from itertools import accumulate
from typing import List, Tuple

import sys


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def getRange(
    nums: List[int],
    *,
    isMax=False,
    isLeftStrict=True,
    isRightStrict=False,
) -> List[Tuple[int, int]]:
    """
    求每个元素作为最值的影响范围(区间)
    默认为每个数作为左严格右非严格最小值的影响区间 [left,right]

    有时候为了避免重复计算，我们可以考虑左侧`严格小于`当前元素的最近元素位置，
    以及右侧`小于等于`当前元素的最近元素位置。
    """

    def compareLeft(stackValue: int, curValue: int) -> bool:
        if isLeftStrict and isMax:
            return stackValue <= curValue
        elif isLeftStrict and not isMax:
            return stackValue >= curValue
        elif not isLeftStrict and isMax:
            return stackValue < curValue
        else:
            return stackValue > curValue

    def compareRight(stackValue: int, curValue: int) -> bool:
        if isRightStrict and isMax:
            return stackValue <= curValue
        elif isRightStrict and not isMax:
            return stackValue >= curValue
        elif not isRightStrict and isMax:
            return stackValue < curValue
        else:
            return stackValue > curValue

    n = len(nums)
    leftMost = [0] * n
    rightMost = [n - 1] * n

    stack = []
    for i in range(n):
        while stack and compareRight(nums[stack[-1]], nums[i]):
            rightMost[stack.pop()] = i - 1
        stack.append(i)

    stack = []
    for i in range(n - 1, -1, -1):
        while stack and compareLeft(nums[stack[-1]], nums[i]):
            leftMost[stack.pop()] = i + 1
        stack.append(i)

    return list(zip(leftMost, rightMost))


n = int(input())
nums = list(map(int, input().split()))
ranges = getRange(nums, isMax=True)


# 每个数产生的不同个数由之前转移过来
dp = [0] * (n + 1)
dp[0] = 1
for i in range(n):
    ndp, dpSum = dp[:], [0] + list(accumulate(dp, lambda x, y: (x + y) % MOD))
    left, right = ranges[i]
    for j in range(left, right + 1):
        ndp[j + 1] += dpSum[j + 1] - dpSum[left]  # j+1里+1 因为前面有个虚拟的0
        ndp[j + 1] %= MOD
    dp = ndp

print(dp[n])
