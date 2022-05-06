# def maxLenAfterRemove(nums: List[int]) -> int:
#     """给出一个数组，最多删除一个连续子数组，求剩下数组的严格递增连续子数组的最大长度。

#     n<=1e6

#     解法2：
#     对位置为i的元素，最大长度为在它之前的所有比它小的元素的值的pre[j]的最大值 + suf[i] nlogn
#     用线段树/树状数组维护小于该元素的最大的pre[j]
#     """
#     n = len(nums)
#     pre = [1] * n
#     suf = [1] * n
#     for i in range(1, n):
#         if nums[i] > nums[i - 1]:
#             pre[i] = pre[i - 1] + 1

#     for i in range(n - 2, -1, -1):
#         if nums[i] < nums[i + 1]:
#             suf[i] = suf[i + 1] + 1


# print(maxLenAfterRemove([5, 3, 4, 9, 2, 8, 6, 7, 1]))  # 4


from collections import defaultdict
from typing import List


class BIT:
    """单点修改"""

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


def maxLenAfterRemove(nums: List[int]) -> int:
    """给出一个数组，最多删除一个连续子数组，求剩下数组的严格递增连续子数组的最大长度。
    
    n<=1e6

    解法 dp + 树状数组：
    dp[i][0] : 以 nums[i] 结尾未删除最大长度
    dp[i][1] : 以 nums[i] 结尾发生删除最大长度

    if nums[i] > nums[i-1] : 
      dp[i][0] = dp[i-1][0] + 1
      dp[i][1] = max(dp[i-1][1] + 1, get(x-1) + 1)
    else : 
        dp[i][0] = 1
        dp[i][1] = get(x-1) + 1
    """
    ...


print(maxLenAfterRemove([5, 3, 4, 9, 2, 8, 6, 7, 1]))  # 4
