from functools import reduce
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的整数数组 nums 和一个正整数 k 。

# 你可以对数组执行以下操作 任意次 ：

# 选择数组里的 任意 一个元素，并将它的 二进制 表示 翻转 一个数位，翻转数位表示将 0 变成 1 或者将 1 变成 0 。
# 你的目标是让数组里 所有 元素的按位异或和得到 k ，请你返回达成这一目标的 最少 操作次数。


# 注意，你也可以将一个数的前导 0 翻转。比方说，数字 (101)2 翻转第四个数位，得到 (1101)2
class Solution:
    def minOperations(self, nums: List[int], k: int) -> int:
        xor_ = reduce(lambda x, y: x ^ y, nums)
        return (xor_ ^ k).bit_count()
