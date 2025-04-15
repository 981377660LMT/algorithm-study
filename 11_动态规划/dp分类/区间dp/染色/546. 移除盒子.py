# 546. 移除盒子
# https://leetcode.cn/problems/remove-boxes/description/
#
# 给出一些不同颜色的盒子 boxes ，盒子的颜色由不同的正数表示。
# 你将经过若干轮操作去去掉盒子，直到所有的盒子都去掉为止。每一轮你可以移除具有相同颜色的连续 k 个盒子（k >= 1），这样一轮之后你将得到 k * k 个积分。
# 返回 你能获得的最大积分和 。
#
# n<=100
# 664. 奇怪的打印机.py

from itertools import groupby
from typing import List
from functools import lru_cache


class Solution:
    def removeBoxes(self, boxes: List[int]) -> int:
        values, counts = [], []
        for v, g in groupby(boxes):
            values.append(v)
            counts.append(len(list(g)))

        @lru_cache(None)
        def dfs(start: int, end: int, k: int) -> int:
            """消除[start:end]中的所有盒子, 且此时和左侧盒子颜色相同的盒子有k个时的最大积分和."""
            if start >= end:
                return 0
            k += counts[start]
            res = k * k + dfs(start + 1, end, 0)  # !直接移除这一段
            for mid in range(start + 1, end):
                if values[mid] == values[start]:
                    res = max(
                        res, dfs(start + 1, mid, 0) + dfs(mid, end, k)
                    )  # !处理这一段后，start和mid合并
            return res

        return dfs(0, len(values), 0)


if __name__ == "__main__":
    print(Solution().removeBoxes(boxes=[1, 3, 2, 2, 2, 3, 4, 3, 1]))
    # 输出：23
    # 解释：
    # [1, 3, 2, 2, 2, 3, 4, 3, 1]
    # ----> [1, 3, 3, 4, 3, 1] (3*3=9 分)
    # ----> [1, 3, 3, 3, 1] (1*1=1 分)
    # ----> [1, 1] (3*3=9 分)
    # ----> [] (2*2=4 分)
