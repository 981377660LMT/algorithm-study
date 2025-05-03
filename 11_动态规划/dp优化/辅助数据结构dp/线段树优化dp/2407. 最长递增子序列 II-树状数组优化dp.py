# https://leetcode.cn/problems/longest-increasing-subsequence-ii/solutions/2022560/shu-zhuang-shu-zu-onlogk-by-ling-jian-20-iqzq/
# !求严格递增的LIS长度 子序列子序列中相邻元素的差值 不超过 k 。
# 1 <= nums.length <= 1e5
# 1 <= nums[i], k <= 1e5
#
# 树状数组O(nlogk)
# !考虑到永远查询的区间都是[i - k, i)形式的，可以按可能的取值范围每k个元素划分一组，也就是[0, k), [k, 2k), [2k, 3k)...
# !对于任意一个[i - k, i)区间的查询，设i = idx * k + r，
# !则它可以分解成[i - k, idx * k)和[idx * k, i)两个区间的查询，
# 如果对于划分过的每一组[nk, (n+k)k)的区间都使用树状数组维护各个前缀和后缀范围的最大值，
# 那么就可以在O(logk)的时间内更新和查询到结果。

from typing import List


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Fenwick:
    __slots__ = ("tree",)

    def __init__(self, n):
        self.tree = [0] * (n + 1)

    def add(self, x, value):
        while x < len(self.tree):
            self.tree[x] = max2(self.tree[x], value)
            x += x & -x

    def get(self, x):
        res = 0
        while x > 0:
            res = max2(res, self.tree[x])
            x -= x & -x
        return res


class Solution:
    def lengthOfLIS(self, nums: List[int], k: int) -> int:
        max_ = max(nums)
        dp = [[Fenwick(k), Fenwick(k)] for _ in range(max_ // k + 1)]
        res = 0
        for i in nums:
            idx, r = divmod(i, k)
            # [idx * k, i)
            a2 = dp[idx][0].get(r)
            # [i - k, idx * k)
            if idx > 0:
                a2 = max(a2, dp[idx - 1][1].get(k - r))
            a2 += 1
            res = max(res, a2)
            dp[idx][0].add(r + 1, a2)
            dp[idx][1].add(k - r, a2)
        return res
