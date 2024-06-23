from itertools import accumulate, permutations
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个整数 n 和一个二维数组 requirements ，其中 requirements[i] = [endi, cnti] 表示这个要求中的末尾下标和 逆序对 的数目。

# 整数数组 nums 中一个下标对 (i, j) 如果满足以下条件，那么它们被称为一个 逆序对 ：

# i < j 且 nums[i] > nums[j]
# 请你返回 [0, 1, 2, ..., n - 1] 的 排列 perm 的数目，满足对 所有 的 requirements[i] 都有 perm[0..endi] 恰好有 cnti 个逆序对。


# 由于答案可能会很大，将它对 109 + 7 取余 后返回。


# !dp[i][j][k] 当前分配到第i个数，使用的


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def numberOfPermutations(self, n: int, requirements: List[List[int]]) -> int:
        requirements.sort()
        for a, b in zip(requirements, requirements[1:]):
            if a[1] > b[1]:
                return 0
        if requirements[0][0] == 0 and requirements[0][1] != 0:
            return 0

        k = requirements[-1][1]
        if k == 0:
            return 1

        requirements.pop()

        limits = {end: cnt for end, cnt in requirements}
        dp = [0] * (k + 1)
        dp[0] = 1
        for i in range(1, n):
            ndp = [0] * (k + 1)
            ndp[0] = 1
            dpSum = [0] + list(accumulate(dp))
            for j in range(1, k + 1):
                ndp[j] = (dpSum[j + 1] - dpSum[max2(0, j - i)]) % MOD
            dp = ndp
            if i in limits:
                cnt = limits[i]
                # print(dp, i, cnt)
                for j in range(k + 1):
                    if j != cnt:
                        dp[j] = 0

        return dp[k - 1] % MOD


# n = 3, requirements = [[2,2],[1,1],[0,0]]

print(Solution().numberOfPermutations(n=3, requirements=[[2, 2], [1, 1], [0, 0]]))
# 5
# [[0,0],[4,3],[1,1]]
print(Solution().numberOfPermutations(n=5, requirements=[[0, 0], [4, 3], [1, 1]]))


def countInversePairs(nums: List[int]) -> int:
    n = len(nums)
    res = 0
    for i in range(n):
        for j in range(i + 1, n):
            if nums[i] > nums[j]:
                res += 1
    return res


def bruteForce(n: int, requirements: List[List[int]]) -> int:
    requirements = sorted(requirements)
    res = 0
    for perm in permutations(range(n)):
        for end, cnt in requirements:
            if countInversePairs(perm[: end + 1]) != cnt:
                break
        else:
            res += 1
    return res


for _ in range(100):
    import random

    n = random.randint(1, 5)
    requirements = [[random.randint(0, n - 1), random.randint(0, 5)] for _ in range(n)]
    res1 = Solution().numberOfPermutations(n, requirements)
    res2 = bruteForce(n, requirements)
    if res1 != res2:
        print(n, requirements)
        print(res1, res2)
        raise
        break

print(bruteForce(n=5, requirements=[[0, 0], [4, 3], [1, 1]]))
