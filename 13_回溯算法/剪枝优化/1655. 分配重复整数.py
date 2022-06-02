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

    def canDistribute2(self, nums: List[int], quantity: List[int]) -> bool:
        """记忆化dp来优化回溯解法
        
        `1815. 得到新鲜甜甜圈的最多组数-回溯+元组记忆化`

        注意这道题比1815多了一个贪心的思想 
        即确定quantity中的每一个数值应该放到哪个容器中 肯定是贪心选择放多的容器 直接把容器数量降成了m

        """

        @lru_cache(None)
        def dfs(index: int, remain: Tuple[int, ...]) -> bool:
            if index == n:
                return True

            for i, num in enumerate(remain):
                if num >= quantity[index]:
                    # !注意这里要保持顺序 否则就起不到记忆化的效果
                    nextRemain = sorted(remain[:i] + (num - quantity[index],) + remain[i + 1 :])
                    if dfs(index + 1, tuple(nextRemain)):
                        return True

            return False

        n = len(quantity)
        remain = sorted(Counter(nums).values(), reverse=True)[:n]

        res = dfs(0, tuple(remain))
        dfs.cache_clear()
        return res


print(Solution().canDistribute2([1, 1, 2, 2], quantity=[2, 2]))
# 输出：true
# 解释：第 0 位顾客得到 [1,1] ，第 1 位顾客得到 [2,2] 。
