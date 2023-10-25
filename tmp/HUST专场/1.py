from bisect import bisect_left
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个非负整数数组 chopsticks 表示筷子长度。如果存在一个数 x ，使得 chopsticks 中恰好有 x 根筷子的长度 大于或者等于 x ，那么就称 chopsticks 是一组 特殊筷子 ，而 x 是该筷子组的 特征值 。


# 注： x 不必 是 chopsticks 的中的元素。
# 如果数组 chopsticks 是一组 特殊筷子 ，请返回它的特征值 x 。否则，返回 -1 。可以证明的是，如果 chopsticks 是一组特殊筷子，那么其特征值 x 是 唯一的 。class Solution:
class Solution:
    def specialChopsticks(self, chopsticks: List[int]) -> int:
        sortedNums = sorted(chopsticks)
        res = -1
        for x in range(len(sortedNums) + 1):
            bigger = len(sortedNums) - bisect_left(sortedNums, x)
            if bigger == x:
                res = x
                break
        return res
