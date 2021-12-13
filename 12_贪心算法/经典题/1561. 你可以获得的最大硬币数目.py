from typing import List

# 每一轮中，你将会选出 任意 3 堆硬币（不一定连续）。
# Alice 将会取走硬币数量最多的那一堆。
# 你将会取走硬币数量第二多的那一堆。
# Bob 将会取走最后一堆。
# 重复这个过程，直到没有更多硬币。
# 给你一个整数数组 piles ，其中 piles[i] 是第 i 堆中硬币的数目。

# 返回你可以获得的最大硬币数目。

# 总结：两大一小
class Solution:
    def maxCoins(self, piles: List[int]) -> int:
        n = len(piles)
        return sum(sorted(piles)[n // 3 :: 2])


print(Solution().maxCoins(piles=[2, 4, 1, 2, 7, 8]))
# 输出：9
# 解释：选出 (2, 7, 8) ，Alice 取走 8 枚硬币的那堆，你取走 7 枚硬币的那堆，Bob 取走最后一堆。
# 选出 (1, 2, 4) , Alice 取走 4 枚硬币的那堆，你取走 2 枚硬币的那堆，Bob 取走最后一堆。
# 你可以获得的最大硬币数目：7 + 2 = 9.
# 考虑另外一种情况，如果选出的是 (1, 2, 8) 和 (2, 4, 7) ，你就只能得到 2 + 4 = 6 枚硬币，这不是最优解。
