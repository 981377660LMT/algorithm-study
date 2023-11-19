from typing import List, Tuple, Optional
from collections import defaultdict, Counter, deque
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)
# 桌子上有 n 个球，每个球的颜色不是黑色，就是白色。

# 给你一个长度为 n 、下标从 0 开始的二进制字符串 s，其中 1 和 0 分别代表黑色和白色的球。

# 在每一步中，你可以选择两个相邻的球并交换它们。

# 返回「将所有黑色球都移到右侧，所有白色球都移到左侧所需的 最小步数」。


def minAdjacentSwap2(nums1: List[int], nums2: List[int]) -> int:
    """求使两个数组相等的最少邻位交换次数

    映射+求逆序对

    时间复杂度`O(nlogn)`
    """

    def countInversionPair(nums: List[int]) -> int:
        """计算逆序对的个数

        sortedList解法

        时间复杂度`O(nlogn)`
        """
        res = 0
        sl = SortedList()
        for num in reversed(nums):
            pos = sl.bisect_left(num)
            res += pos
            sl.add(num)
        return res

    # 含有重复元素的映射 例如nums [1,3,2,1,4] 表示已经排序的数组  [0,1,2,3,4]
    # 那么nums1 [1,1,3,4,2] 就 映射到 [0,3,1,4,2]
    mapping = defaultdict(deque)
    for index, num in enumerate(nums2):
        mapping[num].append(index)

    for index, num in enumerate(nums1):
        mapped = mapping[num].popleft()
        nums1[index] = mapped

    res = countInversionPair(nums1)

    return res


class Solution:
    def minimumSteps(self, s: str) -> int:
        cur = list(int(v) for v in s)
        zero = cur.count(0)
        target = [0] * zero + [1] * (len(cur) - zero)
        return minAdjacentSwap2(cur, target)
