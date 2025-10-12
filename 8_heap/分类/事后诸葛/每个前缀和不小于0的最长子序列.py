# https://leetcode.cn/problems/maximum-transactions-without-negative-balance/?envType=problem-list-v2&envId=sZVESpvF
# 每个前缀和不小于0的最长子序列
# 银行账户
# n<=2000 -1e9<=nums[i]<=1e9

# !堆+反悔 把负数存起来 nlogn(n)
# !也可以dp[i][j] 表示考虑到第i个数，且已经选择了j个数，目前的最大和。 O(n^2)

from heapq import heappop, heappush
from typing import List


def solve(transaction: List[int]) -> int:
    """求每个前缀和都不小于0的最长子序列的长度

    堆+反悔
    """
    res, curSum, pq = 0, 0, []
    for num in transaction:
        res += 1
        curSum += num
        if num < 0:
            heappush(pq, num)
        while curSum < 0 and pq:
            curSum -= heappop(pq)
            res -= 1

    return res


if __name__ == "__main__":
    assert solve([1, 2, 3, -3, -2]) == 5
    assert solve([3, 2, -5, -6, -1, 4]) == 4
