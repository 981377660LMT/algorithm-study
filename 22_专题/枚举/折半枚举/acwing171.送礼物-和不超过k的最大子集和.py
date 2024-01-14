# 达达帮翰翰给女生送礼物，翰翰一共准备了 N 个礼物，其中第 i 个礼物的重量是 G[i]。
# 达达的力气很大，他一次可以搬动重量之和不超过 W 的任意多个物品。
# 达达希望一次搬掉尽量重的一些物品，请你告诉达达在他的力气范围内一次性能搬动的最大重量是多少。

# 1≤N≤46,
# 如果直接枚举的话, 最多有2^46种(如果不考虑W的话)
# 如果对半枚举, 每一半最多2^23=8388608. 可以接受.
# 分成了left和right集合.
# 将right进行排序.
# 对于left的每个元素, 用二分搜索找到right中可以配对的最大元素.

# !折半枚举/折半搜索

from bisect import bisect_right
from typing import List


def solve(weights: List[int], powerLimit: int) -> int:
    def getSortedSubsetSum(nums: List[int], upper: int) -> List[int]:
        res = set([0])
        for num in nums:
            for pre in list(res):
                if pre + num <= upper:
                    res.add(pre + num)
        return sorted(res)

    n = len(weights)
    leftSum = getSortedSubsetSum(weights[: n // 2], powerLimit)
    rightSum = getSortedSubsetSum(weights[n // 2 :], powerLimit)
    res = 0
    for num in leftSum:
        curLimit = powerLimit - num
        pos = bisect_right(rightSum, curLimit) - 1
        if pos >= 0:
            res = max(res, num + rightSum[pos])
    return res


if __name__ == "__main__":
    powerLimit, n = map(int, input().split())
    weights = [int(input()) for _ in range(n)]
    print(solve(weights, powerLimit))
