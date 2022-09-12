# 求最长严格递增子序列(LIS)的个数
# 1 <= nums.length <= 1e5
# -1e6 <= nums[i] <= 1e6

from typing import List
from bisect import bisect_left
from collections import defaultdict


class BIT:
    """单点修改"""

    __slots__ = ("size", "tree")

    def __init__(self, n: int):
        self.size = n
        self.tree = defaultdict(int)

    def add(self, index: int, delta: int) -> None:
        if index <= 0:
            raise ValueError("index 必须是正整数")
        while index <= self.size:
            self.tree[index] += delta
            index += index & -index

    def query(self, index: int) -> int:
        if index > self.size:
            index = self.size
        res = 0
        while index > 0:
            res += self.tree[index]
            index -= index & -index
        return res

    def queryRange(self, left: int, right: int) -> int:
        return self.query(right) - self.query(left - 1)


class Discretizer:
    """离散化"""

    def __init__(self, nums: List[int]) -> None:
        allNums = sorted(set(nums))
        self.mapping = {allNums[i]: i + 1 for i in range(len(allNums))}

    def get(self, num: int) -> int:
        if num not in self.mapping:
            raise ValueError(f"{num} not in {self.mapping}")
        return self.mapping[num]

    def __len__(self) -> int:
        return len(self.mapping)


# O(nlogn)
class Solution:
    def findNumberOfLIS(self, nums: List[int]) -> int:
        """nlogn求最长严格递增子序列的个数"""
        n = len(nums)
        if n <= 1:
            return n

        discretizer = Discretizer(nums)

        LIS = []
        # 每个长度的LIS对应一个BIT，BIT维护结尾小于等于value的子序列有多少个
        dp = defaultdict(lambda: BIT(len(discretizer) + 10))

        for num in nums:
            pos = bisect_left(LIS, num)
            if pos == len(LIS):
                LIS.append(num)
            else:
                LIS[pos] = num

            # 求前一个位置结尾小于当前元素的子序列的个数
            # 遍历可以用树状数组优化
            preBIT = dp[pos - 1]
            count = preBIT.query(discretizer.get(num) - 1)
            dp[pos].add(discretizer.get(num), max(1, count))

        lastPos = len(LIS) - 1
        return dp[lastPos].query(int(1e20))


print(Solution().findNumberOfLIS([1, 3, 2, 5, 4, 7]))
# 输出: 2
# 解释: 有两个最长递增子序列，分别是 [1, 3, 4, 7] 和[1, 3, 5, 7]。
