# https://leetcode.cn/problems/maximum-white-tiles-covered-by-a-carpet/
from BIT import BIT2

from typing import List


class Solution:
    def maximumWhiteTiles(self, tiles: List[List[int]], carpetLen: int) -> int:
        """树状数组区间更新查询"""
        bit = BIT2(int(1e9 + 10))
        res = 0
        for left, right in tiles:
            bit.add(left, right + 1, 1)
        # 毯子的左端点一定和某组瓷砖的左端点一致。
        for left, _ in tiles:
            res = max(res, bit.query(left, left + carpetLen))
        return res
