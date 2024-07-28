# 3224. 使差值相等的最少数组改动次数
# https://leetcode.cn/problems/minimum-array-changes-to-make-differences-equal/description/
# 给你一个长度为 n 的整数数组 nums ，n 是 偶数 ，同时给你一个整数 k 。
# 你可以对数组进行一些操作。每次操作中，你可以将数组中 任一 元素替换为 0 到 k 之间的 任一 整数。
# 执行完所有操作以后，你需要确保最后得到的数组满足以下条件
# 存在一个整数 X ，满足对于所有的 (0 <= i < n) 都有 abs(a[i] - a[n - i - 1]) = X 。
# 请你返回满足以上条件 最少 修改次数。
# !0 <= nums[i] <= k <= 1e5
# !n 是偶数
#
# 1674. 使数组互补的最少操作次数
# https://leetcode.cn/problems/minimum-moves-to-make-array-complementary/solutions/

from collections import defaultdict
from typing import List

MOD = int(1e9 + 7)
INF = int(1e20)


def min2(a: int, b: int) -> int:
    return a if a < b else b


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def minChanges(self, nums: List[int], k: int) -> int:
        diff = defaultdict(int)  # 目标差值为x时，需要的操作数

        def add(left: int, right: int, x: int) -> None:
            diff[left] += x
            diff[right + 1] -= x

        n = len(nums)
        for i in range(n // 2):
            a, b = nums[i], nums[~i]
            if a > b:
                a, b = b, a
            d = b - a
            add(0, d - 1, 1)  # 0 <= x <= d-1 时需要一次操作
            mx = max(a, k - b, b, k - a)
            add(d + 1, mx, 1)
            add(mx + 1, INF, 2)

        res, curSum = INF, 0
        for key in sorted(diff):
            if key > k:
                break
            curSum += diff[key]
            res = min2(res, curSum)
        return res


assert Solution().minChanges(nums=[1, 1, 1, 1], k=2) == 0
