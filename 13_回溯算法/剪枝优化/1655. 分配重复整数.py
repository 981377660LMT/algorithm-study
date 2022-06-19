from functools import lru_cache
from typing import List, Tuple
from collections import Counter, defaultdict


# 1 <= n <= 105
# 1 <= m <= 10  顾客数量只有 10
# quantity[i] 是第 i 位顾客订单的数目
# 请你判断是否能将 nums 中的整数分配给这些顾客，且满足：
# !第 i 位顾客 恰好 有 quantity[i] 个整数。
# !第 i 位顾客拿到的整数都是 相同的 。

# 回溯
# 1. 对customers排序剪枝


class Solution:
    def canDistribute(self, nums: List[int], quantity: List[int]) -> bool:
        """回溯"""

        def bt(index: int) -> bool:
            if index == len(quantity):
                return True
            for freq, count in list(freqCounter.items()):
                if freq >= quantity[index] and count > 0:
                    freqCounter[freq] -= 1
                    freqCounter[freq - quantity[index]] += 1
                    if bt(index + 1):
                        return True
                    freqCounter[freq] += 1
                    freqCounter[freq - quantity[index]] -= 1
            return False

        quantity.sort(reverse=True)
        freqCounter = defaultdict(int, Counter(Counter(nums).values()))
        return bt(0)


print(Solution().canDistribute([1, 1, 2, 2], quantity=[2, 2]))
# 输出：true
# 解释：第 0 位顾客得到 [1,1] ，第 1 位顾客得到 [2,2] 。
