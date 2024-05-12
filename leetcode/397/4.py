from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个数组 nums ，它是 [0, 1, 2, ..., n - 1] 的一个排列 。对于任意一个 [0, 1, 2, ..., n - 1] 的排列 perm ，其 分数 定义为：

# score(perm) = |perm[0] - nums[perm[1]]| + |perm[1] - nums[perm[2]]| + ... + |perm[n - 1] - nums[perm[0]]|


# 返回具有 最低 分数的排列 perm 。如果存在多个满足题意且分数相等的排列，则返回其中字典序最小的一个。
# 按照字典序搜索


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def findPermutation(self, nums: List[int]) -> List[int]:
        @lru_cache(None)
        def dfs(index: int, visited: int, pre: int, first: int) -> int:
            if index == n:
                return abs(pre - nums[first])
            nonlocal next_
            resCost = INF
            for next in range(n):
                if visited & (1 << next):
                    continue
                nextCost = dfs(index + 1, visited | (1 << next), next, first) + abs(
                    pre - nums[next]
                )
                if nextCost < resCost:
                    resCost = nextCost
                    next_[(index, visited, pre)] = (
                        index + 1,
                        visited | (1 << next),
                        next,
                    )
            return resCost

        resCost = INF
        res = [INF]
        n = len(nums)
        for i in range(n):
            next_ = dict()
            tmp = dfs(1, 1 << i, i, i)
            # print(tmp, resCost, first)
            if tmp < resCost:
                resCost = tmp
                curRes = [i]
                curState = (1, 1 << i, i)
                for _ in range(1, n):
                    curState = next_[curState]
                    curRes.append(curState[2])
                res = curRes
        dfs.cache_clear()
        return res


print(Solution().findPermutation(nums=[0, 2, 1]))
# [0,1,2,3]
print(Solution().findPermutation(nums=[0, 1, 2, 3]))
