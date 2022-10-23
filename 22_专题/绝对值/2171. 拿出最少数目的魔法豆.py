"""
给你一个 正 整数数组 beans ，其中每个整数表示一个袋子里装的魔法豆的数目。

请你从每个袋子中 拿出 一些豆子（也可以 不拿出），
使得剩下的 非空 袋子中（即 至少 还有 一颗 魔法豆的袋子）魔法豆的数目 相等 。
一旦魔法豆从袋子中取出，你不能将它放到任何其他的袋子中。

!请你返回你需要拿出魔法豆的 最少数目。

1 <= beans.length <= 1e5
1 <= beans[i] <= 1e5
"""

from itertools import accumulate
from typing import List

INF = int(1e18)


class Solution:
    def minimumRemoval(self, beans: List[int]) -> int:
        """枚举最后和哪个袋子的魔法豆相等
        左边的都要变成0 右边的都要变成nums[i]
        """
        n = len(beans)
        beans.sort()
        preSum = [0] + list(accumulate(beans))
        res = INF
        for i in range(n):
            leftSum = preSum[i]
            rightSum = (preSum[n] - preSum[i + 1]) - (n - i - 1) * beans[i]
            res = min(res, leftSum + rightSum)
        return res

    def minimumRemoval2(self, beans: List[int]) -> int:
        """枚举最后和哪个袋子的魔法豆相等
        拿出最少等于保留最多
        """
        n = len(beans)
        beans.sort()
        maxSave = 0
        for i in range(n):
            maxSave = max(maxSave, beans[i] * (n - i))
        return sum(beans) - maxSave


print(Solution().minimumRemoval(beans=[4, 1, 6, 5]))
