from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 一个长度为 n 下标从 0 开始的整数数组 arr 的 不平衡数字 定义为，在 sarr = sorted(arr) 数组中，满足以下条件的下标数目：

# 0 <= i < n - 1 ，和
# sarr[i+1] - sarr[i] > 1
# 这里，sorted(arr) 表示将数组 arr 排序后得到的数组。

# 给你一个下标从 0 开始的整数数组 nums ，请你返回它所有 子数组 的 不平衡数字 之和。


# 子数组指的是一个数组中连续一段 非空 的元素序列。
class Solution:
    def a(self, s: str) -> List[str]:
        ...
