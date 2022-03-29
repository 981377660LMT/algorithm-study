from typing import List
from collections import defaultdict

MOD = int(1e9 + 7)

# tree直接用dict 省去离散化步骤


class BIT:
    """单点修改的树状数组"""

    def __init__(self, n: int):
        self.size = n
        self.tree = defaultdict(int)

    @staticmethod
    def _lowbit(index: int) -> int:
        return index & -index

    def add(self, index: int, delta: int) -> None:
        if index <= 0:
            raise ValueError('index 必须是正整数')
        while index <= self.size:
            self.tree[index] += delta
            index += self._lowbit(index)

    def query(self, index: int) -> int:
        if index > self.size:
            index = self.size
        res = 0
        while index > 0:
            res += self.tree[index]
            index -= self._lowbit(index)
        return res

    def sumRange(self, left: int, right: int) -> int:
        return self.query(right) - self.query(left - 1)


class Solution:
    def goodTriplets(self, nums: List[int]) -> int:
        """ 
        求nums[a] < nums[b] and nums[c] > nums[d] (a,b,c,d) 对数
        nums is a permutation of the integers 1...n
        """

        n = len(nums)
        leftSmaller = [0] * n
        rightSmaller = [0] * n

        bit1 = BIT(n + 10)
        for i, num in enumerate(nums):
            smaller = bit1.query(num)
            leftSmaller[i] = smaller
            bit1.add(num, 1)

        bit2 = BIT(n + 10)
        for i, num in reversed(list(enumerate(nums))):
            smaller = bit2.query(num)
            rightSmaller[i] = smaller
            bit2.add(num, 1)

        res = 0
        for i in range(1, n - 1):
            leftSmaller[i] += leftSmaller[i - 1]
            # 前半用前缀和，后半枚举分割点
            res += leftSmaller[i] * rightSmaller[i + 1]
            res %= MOD
        return res


print(Solution().goodTriplets(nums=[1, 2, 5, 4, 3]))
