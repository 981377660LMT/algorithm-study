# 1求最长严格递增子序列的个数 LIS个数

# 1 <= nums.length <= 1e5
# -1e6 <= nums[i] <= 1e6


from typing import List
from bisect import bisect_left
from collections import defaultdict


# https://leetcode.com/problems/number-of-longest-increasing-subsequence/discuss/1643753/Python-O(nlogn)-solution-w-detailed-explanation-of-how-to-develop-a-binary-search-solution-from-300
# 记录lis的同时也记录一个counter数组
# counter[pos][value] 记录了lis[pos][value]结尾的最长上升子序列的个数
# 需要将所有lis[pos-1][k]<lis[pos][value]的counter[pos-1][k]累加到counter[pos][value]
# 答案是counter[-1]所有元素之和

# lis：[1, 2, 4, 7]
# counter：defaultdict(<class 'sortedcontainers.sortedlist.SortedList'>,
# {
# -1: SortedList([]),
# 0: SortedList([(1, 1)]),
# 1: SortedList([(2, 1), (3, 1)]),
# 2: SortedList([(4, 2), (5, 2)]),
# 3: SortedList([(7, 4)])
# }


# O(nlogn)
class Solution:
    def findNumberOfLIS(self, nums: List[int]) -> int:
        n = len(nums)
        if n <= 1:
            return n

        LIS = []
        D = Discretizer(nums)
        # 每个长度的LIS对应一个BIT，BIT维护结尾小于等于value的子序列有多少个
        dp = defaultdict(lambda: BIT(len(D) + 10))

        for num in nums:
            pos = bisect_left(LIS, num)
            if pos == len(LIS):
                LIS.append(num)
            else:
                LIS[pos] = num

            # 上一个位置结尾小于当前元素的所有的子序列的个数是多少
            # 遍历可以用树状数组优化
            preBIT = dp[pos - 1]
            count = preBIT.query(D.get(num) - 1)
            dp[pos].add(D.get(num), count if count > 0 else 1)

        lastPos = len(LIS) - 1
        return dp[lastPos].query(len(D))


class BIT:
    """单点修改"""

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


print(Solution().findNumberOfLIS([1, 3, 5, 4, 7]))
# 输出: 2
# 解释: 有两个最长递增子序列，分别是 [1, 3, 4, 7] 和[1, 3, 5, 7]。
