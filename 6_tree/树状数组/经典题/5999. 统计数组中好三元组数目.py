from typing import List, Tuple
from collections import defaultdict


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


# 三元组一般是`定一移二`的做法
# 或者枚举贡献

# 1. 将nums2数组的值映射到nums1对应的数的索引
# 2. 使用树状数组处理nums2每个数左边有多少个比他小，右边有多少个比他大，得到leftSmaller和rightBigger两个数组
# 3. 枚举每个数作为中间数，有leftSmaller[i]*rightBigger[i]种取法，求和即可
# 时间复杂度为O(nlogn)，空间复杂度为O(n)
# 其中第一步的映射操作有点像 1713. 得到子序列的最少操作次数
class Solution:
    def goodTriplets(self, nums1: List[int], nums2: List[int]) -> int:
        n = len(nums1)
        indexByValue = defaultdict(lambda: -1)
        for i, num in enumerate(nums1):
            indexByValue[num] = i
        target = [indexByValue[num] for num in nums2]

        leftSmaller = [0] * n
        rightBigger = [0] * n

        bit1 = BIT(n + 10)
        for i, num in enumerate(target):
            smaller = bit1.query(num + 1)
            leftSmaller[i] = smaller
            bit1.add(num + 1, 1)

        bit2 = BIT(n + 10)
        for i in range(n - 1, -1, -1):
            bigger = bit2.sumRange(target[i] + 2, n + 1)
            rightBigger[i] = bigger
            bit2.add(target[i] + 1, 1)

        res = 0
        for left, right in zip(leftSmaller, rightBigger):
            res += left * right

        return res


print(Solution().goodTriplets(nums1=[2, 0, 1, 3], nums2=[0, 1, 2, 3]))
print(Solution().goodTriplets(nums1=[4, 0, 1, 3, 2], nums2=[4, 1, 0, 2, 3]))
