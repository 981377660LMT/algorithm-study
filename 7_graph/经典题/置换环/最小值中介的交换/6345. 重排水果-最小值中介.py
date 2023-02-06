# 你有两个果篮，每个果篮中有 n 个水果。
# 给你两个下标从 0 开始的整数数组 basket1 和 basket2 ，用以表示两个果篮中每个水果的成本。
# 你希望两个果篮相等。为此，可以根据需要多次执行下述操作：
# 选中两个下标 i 和 j ，并交换 basket1 中的第 i 个水果和 basket2 中的第 j 个水果。
# !交换的成本是 min(basket1i,basket2j) 。
# 根据果篮中水果的成本进行排序，如果排序后结果完全相同，则认为两个果篮相等。
# !返回使两个果篮相等的最小交换成本，如果无法使两个果篮相等，则返回 -1 。

# 解法:
# !最小值中介


from typing import List
from collections import Counter


class Solution:
    def minCost(self, basket1: List[int], basket2: List[int]) -> int:
        counter1, counter2 = Counter(basket1), Counter(basket2)
        counter = counter1 + counter2
        if any(v & 1 for v in counter.values()):  # 每种必须是偶数
            return -1
        target = Counter({k: v // 2 for k, v in counter.items()})
        diff1 = sorted((target - counter1).elements())  # 缺少的
        diff2 = sorted((counter1 - target).elements(), reverse=True)  # 多出的
        min_ = min(counter)
        return sum(min(a, b, min_ * 2) for a, b in zip(diff1, diff2))  # 田忌赛马


assert (
    Solution().minCost(
        basket1=[84, 80, 43, 8, 80, 88, 43, 14, 100, 88],
        basket2=[32, 32, 42, 68, 68, 100, 42, 84, 14, 8],
    )
) == 48
